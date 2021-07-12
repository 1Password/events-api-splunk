# 1Password Events API for Splunk Installation Steps

1. Download the 1Password Splunk zip file from the provided link and unzip it with the software of your choosing.
2. Sign into your Splunk account, and click the gear icon next to the apps listing on the left-hand side of the home screen to navigate to App Management.
3. Click 'Install App from File'
4. Select the relevant 1Password Splunk _tar_ file for your OS and click 'Upload'
5. After a successful install, click 'Set up now' to set up the app
6. Fill in the token generated from your 1Password account
7. Click 'Submit'
8. Go to Settings -> Data Inputs -> Scripts -> and click Enable on the 1Password inputs you're interested in. Not that these scripts log to the default index unless configured to log elsewhere. Check out our [support page](https://support.1password.com/cs/events-reporting-splunk/#step-3-set-up-the-1password-events-api-add-on) for more information.
9. Click the 1Password Splunk app in the apps listing and search `sourcetype="1password:insights:signin_attempts"` or `sourcetype="1password:insights:item_usages"` to see your 1Password events. This data is ready to be analyzed, monitored and reported.
