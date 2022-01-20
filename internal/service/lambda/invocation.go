package lambda

import (
	"crypto/md5"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceInvocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceInvocationCreate,
		Read:   resourceInvocationRead,
		Update: resourceInvocationUpdate,
		Delete: resourceInvocationDelete,

		Schema: map[string]*schema.Schema{
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"qualifier": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  FunctionVersionLatest,
			},

			"input": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},

			"triggers": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceInvocationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).LambdaConn

	functionName := d.Get("function_name").(string)
	qualifier := d.Get("qualifier").(string)
	input := []byte(d.Get("input").(string))

	res, err := conn.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: aws.String(lambda.InvocationTypeRequestResponse),
		Payload:        input,
		Qualifier:      aws.String(qualifier),
	})

	if err != nil {
		return err
	}

	if res.FunctionError != nil {
		return fmt.Errorf("Lambda function (%s) returned error: (%s)", functionName, string(res.Payload))
	}

	if err = d.Set("result", string(res.Payload)); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s_%s_%x", functionName, qualifier, md5.Sum(input)))

	return nil
}

func resourceInvocationRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceInvocationUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceInvocationCreate(d, meta)
}

func resourceInvocationDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	d.Set("result", nil)
	return nil
}
