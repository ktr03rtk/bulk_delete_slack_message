# bulk_delete_slack_message

## Overview

This program delete SLACK messages older than the entered date.

## Usage

- Prepare the SLACK token with correct scope. This program uses conversations.info
  , conversations.history, chat.delete API. To attach correct scope, see the official document.

- Definite the slack.env file by referring to the sample_slack.env file.

- Run the `docker-compose run bulk-delete` command to start the program. It inquire the latest timestamp for deletion, enter the timestamp and confirm.
