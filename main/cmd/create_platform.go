package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

var verbose bool
var profile string
var dryRun bool
var chartName string
var chartVersion string
var extraRequirements string

var ingressIP string
var ingressNodeSelector string

var caCertEmail string

func init() {
	createPlatformCmd.Flags().BoolVarP(&dryRun, "verbose", "v", false, "Indicates if a verbose output mode should be used.")
	createPlatformCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Indicates if a dry run should be used i.e. kif should generate Helm charts without executing them.")
	createPlatformCmd.Flags().StringVar(&profile, "cloud", "baremetal", "Cloud provider to use.")
	createPlatformCmd.Flags().StringVar(&chartName, "chart-name", "kif", "Name of the generated chart.")
	createPlatformCmd.Flags().StringVar(&chartVersion, "chart-version", "0.0.0", "Version of the generated chart.")
	createPlatformCmd.Flags().StringVar(&extraRequirements, "extra-requirements", "", "Extra requirements to be included in generated chart requirements file.")

	createPlatformCmd.Flags().StringVar(&ingressIP, "ingress-ip", "", "IP address of ingress node.")
	createPlatformCmd.Flags().StringVar(&ingressNodeSelector, "ingress-node-selector", "machine0001", "Node selector of ingress pod.")

	createPlatformCmd.Flags().StringVar(&caCertEmail, "cert-email", "", "CA certificate administrator e-mail used during ACME registration process.")

	rootCmd.AddCommand(createPlatformCmd)
}

var createPlatformCmd = &cobra.Command{
	Use:   "create platform",
	Short: "Create kif platform.",
	Long:  `Create kif platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		kif := OrExitOnError(NewKifPlatform()).(*KifPlatform)

		if ingressIP == "" {
			fmt.Println("Ingress IP cannot be empty. Please use --ingress-ip option.")
			os.Exit(-1)
		}
		if caCertEmail == "" {
			fmt.Println("CA cert admin e-mail cannot be empty. Please use --cert-email option.")
			os.Exit(-1)
		}
		kif.Configuration = map[string]interface{}{
			"Chart": map[string]interface{}{
				"Name":    chartName,
				"Version": chartVersion,
			},
			"Ingress": map[string]interface{}{
				"ExternalIp":   ingressIP,
				"NodeSelector": ingressNodeSelector,
			},
			"CertManager": map[string]interface{}{
				"Email": caCertEmail,
			},
			"Prometheus": map[string]interface{}{
				"Host": fmt.Sprintf("prometheus.%s.nip.io", ingressIP),
			},
		}

		ExitOnError(kif.RenderTemplate("Chart.yaml"))
		ExitOnError(kif.RenderRequirements(extraRequirements))

		valuesTemplateText, err := kif.TemplatesBox.String("values.yml")
		if err != nil {
			fmt.Println(err)
			return
		}
		valuesTemplate, err := template.New("valuesTemplate").Parse(valuesTemplateText)
		if err != nil {
			fmt.Println(err)
			return
		}
		valuesFile, err := os.Create(kif.Sandbox + "/values.yml")
		if err != nil {
			fmt.Println(err)
			return
		}

		ExitMessageOnError(
			kif.RenderTemplate("templates/issuer-letsencrypt.yml"),
			"Cannot generate Let's Encrypt ACME issuer")

		err = valuesTemplate.Execute(valuesFile, kif.Configuration)
		ExitOnError(err)
		command := exec.Command("htpasswd", "-c", "-b", kif.Sandbox+"/auth", "admin", "admin")
		commandOutput, err := command.CombinedOutput()
		ExitOnError(err)
		if verbose {
			println("Generating basic auth authentication for Prometheus:")
			println(string(commandOutput))
		}
		prometheusAuthSecretTemplateFile, err := kif.TemplatesBox.String("secret-ingress-auth-prometheus.yml")
		ExitOnError(err)
		prometheusIngressAuthTemplate, err := template.New("prometheusAuthSecretTemplate").Parse(prometheusAuthSecretTemplateFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		prometheusIngressAuthFile, err := os.Create(kif.Sandbox + "/templates/secret-ingress-auth-prometheus.yml")
		ExitOnError(err)
		auth := OrExitOnError(ioutil.ReadFile(kif.Sandbox + "/auth")).([]byte)
		kif.Configuration["Prometheus"] = map[string]interface{}{
			"Ingress": map[string]interface{}{
				"Auth": base64.StdEncoding.EncodeToString(auth),
			},
		}
		err = prometheusIngressAuthTemplate.Execute(prometheusIngressAuthFile, kif.Configuration)
		ExitOnError(err)

		if dryRun {
			println("Platform chart generated successfully: " + kif.Sandbox)
		} else {
			command = exec.Command("helm", "dependency", "update", kif.Sandbox)
			commandOutput, _ = command.CombinedOutput()
			println(string(commandOutput))

			command = exec.Command("helm", "install", "--namespace=kube-system", "--name=skrt", kif.Sandbox, "--values="+kif.Sandbox+"/values.yml")
			commandOutput, _ = command.CombinedOutput()
			println(string(commandOutput))
		}
	},
}

// Kif platform

type KifPlatform struct {
	Sandbox       string
	TemplatesBox  *rice.Box
	Configuration map[string]interface{}
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
		Sandbox:       sandbox,
		TemplatesBox:  templateBox,
		Configuration: map[string]interface{}{},
	}, nil
}

func (kif *KifPlatform) RenderTemplate(name string) error {
	templateText, err := kif.TemplatesBox.String(name)
	if err != nil {
		return err
	}
	parsedTemplate, err := template.New(name).Parse(templateText)
	if err != nil {
		return err
	}
	targetFile, err := os.Create(kif.Sandbox + "/" + name)
	if err != nil {
		return err
	}
	err = parsedTemplate.Execute(targetFile, kif.Configuration)
	if err != nil {
		return err
	}
	return nil
}

func (kif *KifPlatform) RenderRequirements(extraRequirements string) error {
	requirementsYaml, err := kif.TemplatesBox.Bytes("requirements.yaml")
	if err != nil {
		return err
	}
	requirements := map[string]interface{}{}
	err = yaml.Unmarshal(requirementsYaml, &requirements)
	if err != nil {
		return err
	}
	if extraRequirements != "" {
		extraRequirementsYaml, err := ioutil.ReadFile(extraRequirements)
		if err != nil {
			return err
		}
		extraRequirementsMap := map[string]interface{}{}
		err = yaml.Unmarshal(extraRequirementsYaml, &extraRequirementsMap)
		if err != nil {
			return err
		}
		for _, dependency := range extraRequirementsMap["dependencies"].([]interface{}) {
			requirements["dependencies"] = append(requirements["dependencies"].([]interface{}), dependency)
		}
	}
	requirementsYaml, err = yaml.Marshal(requirements)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(kif.Sandbox+"/requirements.yaml", requirementsYaml, 0644)
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

func ExitMessageOnError(err error, message string) {
	if err != nil {
		fmt.Println(fmt.Sprintf("%s: %s", message, err))
		os.Exit(-1)
	}
}

func OrExitOnError(value interface{}, err error) interface{} {
	ExitOnError(err)
	return value
}
