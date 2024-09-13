# Creating iCalendar Files

**ical-git**'s simple design means it does not directly create iCalendar files. Instead, it relies on external sources to generate these files. Common sources include:

- **Emails**: Many email clients and services allow you to create calendar events directly from email invitations. These events can be exported as iCalendar files.
- **Appointments**: Calendar applications like Google Calendar, Outlook, or Apple Calendar can create and export events in the iCalendar format.
- **Language Models (LLMs)**: Private LLMs can be used to generate iCalendar content based on prompts or templates. Ensure the generated content adheres to the iCalendar format, including proper syntax and required fields.

To use these sources:

1. **Export from Source**: Export the calendar event or generate the content from the LLM in the iCalendar format (`.ics` file).
2. **Place in Directory**: Save the generated or exported `.ics` file in the directory specified in your `icalgit.toml` configuration file.

Alternatively, you can copy and modify existing iCalendar files from the `testdata` directory to suit your needs:

1. **Copy Files**: Copy the desired iCalendar files from the `testdata` directory to your working directory.
2. **Modify Content**: Open the copied files and modify the event details, dates, and other relevant information to fit your requirements.
3. **Save Changes**: Save the modified files in the directory specified in your configuration.

# Managing iCal Files

For detailed instructions on managing iCal files, refer to the [Managing iCal Files](ical.md#managing-ical-files) section.

It is highly advisable to place your iCal files under revision control to ensure that changes are tracked and can be reverted if necessary. 
Additionally, setting up a cron job to periodically pull the latest content of these files can help keep your calendar up-to-date.

## Setting Up Revision Control

1. **Initialize a Git Repository**: Navigate to your iCal files directory and initialize a Git repository.

    ```console
    cd /home/icalgit/icalgit/my-ical-files
    git init
    ```

2. **Add and Commit Your Files**: Add your iCal files to the repository and commit them.

    ```console
    git add .
    git commit -m "Initial commit of iCal files"
    ```

3. **Push to a Remote Repository**: If you have a remote repository (e.g., on GitHub), push your local repository to the remote.

    ```console
    git remote add origin git@github.com:yourusername/your-repo.git
    git push -u origin main
    ```

## Setting Up a Cron Job

To ensure your iCal files are periodically updated, you can set up a cron job to pull the latest content from your revision control system.

1. **Edit Your Crontab**: Open your crontab file for editing.

    ```console
    crontab -e -u icalgit
    ```

2. **Add the Cron Job**: Add a cron job to pull the latest changes from your repository. It's advisable to use a passphraseless SSH key to avoid issues with cron executing in a non-interactive session.

    ```cron
    * * * * * cd /home/icalgit/icalgit/my-ical-files && GIT_SSH_COMMAND="ssh -i /home/icalgit/.ssh/id_icalgit_nopassphrase -o IdentitiesOnly=yes" git pull origin main
    ```

This cron job will run every minute, pulling the latest changes from the `main` branch of your repository. Adjust the schedule as needed.