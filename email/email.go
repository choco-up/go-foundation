package email

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

const CHARSET = "UTF-8"

func VerifyEmail(sesClient *ses.SES, email string) error {
	_, err := sesClient.VerifyEmailIdentity(&ses.VerifyEmailIdentityInput{
		EmailAddress: aws.String(email),
	})

	return err
}

func GetAllSESTemplate(sesClient *ses.SES, page int64) ([]*ses.TemplateMetadata, error) {
	/*
	  Pagination Limit Set 10
	*/
	var itemsFrom int64
	itemsFrom = page*10 - 9
	listTemplatesInput := ses.ListTemplatesInput{
		MaxItems: &itemsFrom,
	}

	listTemplatesOutput, err := sesClient.ListTemplates(&listTemplatesInput)
	if err != nil {
		fmt.Println("Error list ses templates: ", err)
		return nil, err
	}
	return listTemplatesOutput.TemplatesMetadata, nil
}

func SendEmail(sesClient *ses.SES, toAddresses []*string, emailText string, sender string, subject string) error {

	_, err := sesClient.SendEmail(&ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: toAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(CHARSET),
					Data:    aws.String(emailText),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CHARSET),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	})

	return err
}

// https://manishpaneri.blogspot.com/2019/07/how-to-use-aws-ses-template-using-golang.html
func SendHTMLEmail(sesClient *ses.SES, toAddresses []*string, htmlText string, sender string, subject string) error {

	_, err := sesClient.SendEmail(&ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: toAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CHARSET),
					Data:    aws.String(htmlText),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CHARSET),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	})

	return err
}

// GetSESTemplateByName ... get template by using template name
func GetSESTemplateByName(sesClient *ses.SES, templateName string) (map[string]string, error) {
	templateOutput, err := sesClient.GetTemplate(&ses.GetTemplateInput{TemplateName: aws.String(templateName)})
	if err != nil {
		fmt.Println("Error get ses template: ", err)
		return nil, err
	}

	template := map[string]string{
		"templateName": *templateOutput.Template.TemplateName,
		"subject":      *templateOutput.Template.SubjectPart,
		"htmlBody":     *templateOutput.Template.HtmlPart,
		"textBody":     *templateOutput.Template.TextPart,
	}
	return template, nil
}

func SendEmailTest(sesClient *ses.SES, sender string, receivers []string) error {
	var toAddresses []*string
	for _, receiver := range receivers {
		toAddresses = append(toAddresses, aws.String(receiver))
	}
	emailText := "Hello World"
	//sender := "abhishek@learnaws.org"
	subject := "Testing Email!"
	err := SendEmail(sesClient, toAddresses, emailText, sender, subject)
	if err != nil {
		fmt.Printf("Got an error while trying to send email: %v", err)
		return err
	}
	return nil
}

func SendHTMLEmailTest(sesClient *ses.SES, sender string, receivers []string) error {
	var toAddresses []*string
	for _, receiver := range receivers {
		toAddresses = append(toAddresses, aws.String(receiver))
	}
	htmlText := `<html>
	<head></head>
	<h1 style='text-align:center'>This is the heading</h1>
	<p>Hello, world</p>
	</body>
	</html>`
	subject := "Testing Email!"
	err := SendHTMLEmail(sesClient, toAddresses, htmlText, sender, subject)
	if err != nil {
		fmt.Printf("Got an error while trying to send email: %v", err)
		return err
	}
	return nil
}

func SendTemplateEmail(sesClient *ses.SES,
	sender string,
	toAddresses []string,
	ccAddresses []string,
	bccAddresses []string,
	templateName, data string) error {

	var toAddresses_ []*string
	var ccAddresses_ []*string
	var bccAddresses_ []*string

	for _, address := range toAddresses {
		toAddresses_ = append(toAddresses_, &address)
	}

	for _, address := range ccAddresses {
		ccAddresses_ = append(ccAddresses_, &address)
	}

	for _, address := range bccAddresses {
		bccAddresses_ = append(bccAddresses_, &address)
	}

	input := &ses.SendTemplatedEmailInput{
		ConfigurationSetName: nil,
		Destination: &ses.Destination{
			//BccAddresses: []*string{},
			//CcAddresses:  []*string{},
			BccAddresses: bccAddresses_,
			CcAddresses:  ccAddresses_,
			ToAddresses:  toAddresses_,
		},
		Source:       aws.String(sender),
		Template:     &templateName,
		TemplateData: &data,
	}

	_, err := sesClient.SendTemplatedEmail(input)
	if err != nil {
		fmt.Printf("Got an error while trying to send templated email: %v", err)
		return err
	}
	return nil
}
