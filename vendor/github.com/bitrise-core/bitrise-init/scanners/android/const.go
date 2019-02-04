package android

const deployWorkflowDescription = `## How to get a signed APK

This workflow contains the **Sign APK** step. To sign your APK all you have to do is to:

1. Click on **Code Signing** tab
1. Find the **ANDROID KEYSTORE FILE** section
1. Click or drop your file on the upload file field
1. Fill the displayed 3 input fields:
 1. **Keystore password**
 1. **Keystore alias**
 1. **Private key password**
1. Click on **[Save metadata]** button

That's it! From now on, **Sign APK** step will receive your uploaded files.

## To run this workflow

If you want to run this workflow manually:

1. Open the app's build list page
2. Click on **[Start/Schedule a Build]** button
3. Select **deploy** in **Workflow** dropdown input
4. Click **[Start Build]** button

Or if you need this workflow to be started by a GIT event:

1. Click on **Triggers** tab
2. Setup your desired event (push/tag/pull) and select **deploy** workflow
3. Click on **[Done]** and then **[Save]** buttons

The next change in your repository that matches any of your trigger map event will start **deploy** workflow.
`
