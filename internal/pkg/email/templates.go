package email

import "html/template"

// emailTemplates menyimpan template HTML dan teks
var emailTemplates = map[string]*template.Template{
	"verification_html": template.Must(template.New("verification_html").Parse(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Verify Your Email</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				line-height: 1.6;
				color: #333;
			}

			.container {
				max-width: 600px;
				margin: 0 auto;
				padding: 20px;
			}

			.header {
				background: #007bff;
				color: white;
				padding: 20px;
				text-align: center;
				border-radius: 5px 5px 0 0;
			}

			.content {
				background: #f9f9f9;
				padding: 30px;
				border-radius: 0 0 5px 5px;
			}

			.button {
				display: inline-block;
				background: #28a745;
				color: white;
				padding: 12px 24px;
				text-decoration: none;
				border-radius: 5px;
				margin: 20px 0;
			}

			.footer {
				text-align: center;
				margin-top: 20px;
				font-size: 12px;
				color: #666;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>{{.AppName}}</h1>
			</div>
			<div class="content">
				<h2>Hi {{.FirstName}}!</h2>
				<p>Welcome to {{.AppName}}! Please verify your email address by clicking the button below:</p>

				<a href="{{.VerificationURL}}" class="button">Verify My Email</a>

				<p>Or copy and paste this link in your browser:</p>
				<p style="word-break: break-all; background: #eee; padding: 10px; border-radius: 3px;">{{.VerificationURL}}
				</p>

				<p><strong>This link expires in {{.ExpiresIn}}.</strong></p>

				<p>If you didn't create an account with us, please ignore this email.</p>

				<p>Best regards,<br>The {{.AppName}} Team</p>
			</div>
			<div class="footer">
				<p>Need help? <a href="{{.SupportURL}}">Contact Support</a></p>
				<p>{{.AppName}} - {{.AppURL}}</p>
			</div>
		</div>
	</body>
	</html>`)),
	"verification_text": template.Must(template.New("verification_text").Parse(`Hi {{.FirstName}}!Welcome to {{.AppName}}! Please verify your email address by clicking the link below:{{.VerificationURL}}This link expires in {{.ExpiresIn}}.If you didn't create an account with us, please ignore this email.Best regards,The {{.AppName}} TeamNeed help? Contact Support: {{.SupportURL}}{{.AppName}} - {{.AppURL}}`)),
	"welcome_html": template.Must(template.New("welcome_html").Parse(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Welcome!</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				line-height: 1.6;
				color: #333;
			}

			.container {
				max-width: 600px;
				margin: 0 auto;
				padding: 20px;
			}

			.header {
				background: #28a745;
				color: white;
				padding: 20px;
				text-align: center;
				border-radius: 5px 5px 0 0;
			}

			.content {
				background: #f9f9f9;
				padding: 30px;
				border-radius: 0 0 5px 5px;
			}

			.button {
				display: inline-block;
				background: #007bff;
				color: white;
				padding: 12px 24px;
				text-decoration: none;
				border-radius: 5px;
				margin: 20px 0;
			}

			.footer {
				text-align: center;
				margin-top: 20px;
				font-size: 12px;
				color: #666;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Welcome to {{.AppName}}! 🎉</h1>
			</div>
			<div class="content">
				<h2>Hi {{.FirstName}}!</h2>
				<p>Your email has been verified successfully! You can now login to your account and start using
					{{.AppName}}.</p> <a href="{{.AppURL}}" class="button">Login to Your Account</a>
				<p>If you have any questions, feel free to contact our support team.</p>
				<p>Best regards,<br>The {{.AppName}} Team</p>
			</div>
			<div class="footer">
				<p>Need help? <a href="{{.SupportURL}}">Contact Support</a></p>
				<p>{{.AppName}} - {{.AppURL}}</p>
			</div>
		</div>
	</body>
	</html>`)),
	"welcome_text": template.Must(template.New("welcome_text").Parse(`Welcome to {{.AppName}}! 🎉Hi {{.FirstName}}!Your email has been verified successfully! You can now login to your account and start using {{.AppName}}.Login here: {{.AppURL}}If you have any questions, feel free to contact our support team.Best regards,The {{.AppName}} TeamNeed help? Contact Support: {{.SupportURL}}{{.AppName}} - {{.AppURL}}`)),
}

type EmailData struct {
	AppName         string
	FirstName       string
	VerificationURL string
	AppURL          string
	SupportURL      string
	ExpiresIn       string
}