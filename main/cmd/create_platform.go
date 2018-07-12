package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/spf13/cobra"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

var verbose bool
var profile string
var dryRun bool
var ingressIP string
var ingressNodeSelector string

var caCertEmail string

func init() {
	createPlatformCmd.Flags().BoolVarP(&dryRun, "verbose", "v", false, "Indicates if a verbose output mode should be used.")
	createPlatformCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Indicates if a dry run should be used i.e. kif should generate Helm charts without executing them.")
	createPlatformCmd.Flags().StringVar(&profile, "cloud", "baremetal", "Cloud provider to use.")
	createPlatformCmd.Flags().StringVar(&ingressIP, "ingress-ip", "", "IP address of ingress node.")
	createPlatformCmd.Flags().StringVar(&ingressNodeSelector, "ingressNodeSelector", "machine0001", "Node selector of ingress pod.")
	createPlatformCmd.Flags().StringVar(&caCertEmail, "cert-email", "", "CA certificate administrator e-mail used during ACME registration process.")
	rootCmd.AddCommand(createPlatformCmd)
}

var createPlatformCmd = &cobra.Command{
	Use:   "create platform",
	Short: "Create kif platform.",
	Long:  `Create kif platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		kifPlatform := OrExitOnError(NewKifPlatform()).(*KifPlatform)

		chart := OrExitOnError(kifPlatform.TemplatesBox.Bytes("Chart.yaml")).([]byte)
		ExitOnError(ioutil.WriteFile(kifPlatform.Sandbox+"/Chart.yaml", chart, 0644))

		requirements, err := kifPlatform.TemplatesBox.Bytes("requirements.yaml")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = ioutil.WriteFile(kifPlatform.Sandbox+"/requirements.yaml", requirements, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}

		valuesTemplateText, err := kifPlatform.TemplatesBox.String("values.yml")
		if err != nil {
			fmt.Println(err)
			return
		}
		valuesTemplate, err := template.New("valuesTemplate").Parse(valuesTemplateText)
		if err != nil {
			fmt.Println(err)
			return
		}
		valuesFile, err := os.Create(kifPlatform.Sandbox + "/values.yml")
		if err != nil {
			fmt.Println(err)
			return
		}
		if ingressIP == "" {
			fmt.Println("Ingress IP cannot be empty. Please use --ingress-ip option.")
			os.Exit(-1)
		}
		if caCertEmail == "" {
			fmt.Println("CA cert admin e-mail cannot be empty. Please use --cert-email option.")
			os.Exit(-1)
		}

		kifConfiguration := map[string]map[string]interface{}{
			"Ingress": {
				"ExternalIp":   ingressIP,
				"NodeSelector": ingressNodeSelector,
			},
			"CertManager": {
				"Email": caCertEmail,
			},
			"Prometheus": {
				"Host": fmt.Sprintf("prometheus.%s.nip.io", ingressIP),
			},
		}
		err = valuesTemplate.Execute(valuesFile, kifConfiguration)
		ExitOnError(err)

		command := exec.Command("htpasswd", "-c", "-b", kifPlatform.Sandbox+"/auth", "admin", "admin")
		commandOutput, err := command.CombinedOutput()
		ExitOnError(err)
		if verbose {
			println("Generating basic auth authentication for Prometheus:")
			println(string(commandOutput))
		}

		ExitOnError(kifPlatform.RenderTemplate("templates/issuer-letsencrypt", kifConfiguration))

		prometheusAuthSecretTemplateFile, err := kifPlatform.TemplatesBox.String("secret-ingress-auth-prometheus.yml")
		ExitOnError(err)
		prometheusIngressAuthTemplate, err := template.New("prometheusAuthSecretTemplate").Parse(prometheusAuthSecretTemplateFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		prometheusIngressAuthFile, err := os.Create(kifPlatform.Sandbox + "/templates/secret-ingress-auth-prometheus.yml")
		ExitOnError(err)
		auth := OrExitOnError(ioutil.ReadFile(kifPlatform.Sandbox + "/auth")).([]byte)
		kifConfiguration["Prometheus"]["Ingress"] = map[string]interface{}{
			"Auth": base64.StdEncoding.EncodeToString(auth),
		}
		err = prometheusIngressAuthTemplate.Execute(prometheusIngressAuthFile, kifConfiguration)
		ExitOnError(err)

		if dryRun {
			println("Platform chart generated successfully: " + kifPlatform.Sandbox)
		} else {
			command = exec.Command("helm", "dependency", "update", kifPlatform.Sandbox)
			commandOutput, _ = command.CombinedOutput()
			println(string(commandOutput))

			command = exec.Command("helm", "install", "--namespace=kube-system", "--name=skrt", kifPlatform.Sandbox, "--values="+kifPlatform.Sandbox+"/values.yml")
			commandOutput, _ = command.CombinedOutput()
			println(string(commandOutput))
		}
	},
}

// Kif platform

type KifPlatform struct {
	Sandbox      string
	TemplatesBox *rice.Box
}

func NewKifPlatform() (*KifPlatform, error) {
	sandbox := fmt.Sprintf("/tmp/kif_%d", time.Now().Unix())
	err := os.MkdirAll(sandbox+"/templates", 0700)
	if err != nil {
		return nil, err
	}

	templateBox, err := rice.FindBox("templates")
	if err != nil {
		return nil, err
	}

	return &KifPlatform{
		Sandbox:      sandbox,
		TemplatesBox: templateBox,
	}, nil
}

func (kif *KifPlatform) RenderTemplate(name string, kifConfiguration interface{}) error {
	templateText, err := kif.TemplatesBox.String(name + ".yml")
	if err != nil {
		return err
	}
	parsedTemplate, err := template.New(name).Parse(templateText)
	if err != nil {
		return err
	}
	targetFile, err := os.Create(kif.Sandbox + "/" + name + ".yml")
	if err != nil {
		return err
	}
	err = parsedTemplate.Execute(targetFile, kifConfiguration)
	if err != nil {
		return err
	}
	return nil
}

// Helper

func ExitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func OrExitOnError(value interface{}, err error) interface{} {
	ExitOnError(err)
	return value
}
