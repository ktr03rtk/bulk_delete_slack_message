# bulk_delete_slack_message

## Overview

This program delete Slack messages older than the entered date.

## Usage

- Prepare the Slack token with correct scope. This program uses conversations.list, conversations.history, chat.delete API. To attach correct scope, see the Slack official document.

- Add the slack.env file by referring to the sample_slack.env file.

- Run the `docker-compose run message-delete` command to start the program. It inquire the latest timestamp for deletion, enter the timestamp and confirm.
