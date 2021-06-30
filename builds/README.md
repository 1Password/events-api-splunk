# 1Password Events API for Splunk Installation Steps

1. Download the 1Password Splunk zip file from the provided link and unzip it with the software of your choosing.
2. Sign into your Splunk account, and click the gear icon next to the apps listing on the left-hand side of the home screen to navigate to App Management.
3. Click 'Install App from File'
4. Select the relevant 1Password Splunk _tar_ file for your OS and click 'Upload'
5. After a successful install, click 'Set up now' to set up the app
6. Fill in the token generated from the 1Password.com and make sure the base URL is correct.
7. Click 'Submit'
8. You should receive a message confirming changes were made and for changes to take effect, Splunk must be restarted.
9. After signing back into Splunk, click the 1Password Splunk app in the apps listing and search `sourcetype="1password:insights:signin_attempts"` to see 1Password failed sign-in attempt logs. This data is ready to be analyzed, monitored and reported.
