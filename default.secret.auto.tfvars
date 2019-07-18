# Copy me to ./env/aws-*/.secret.auto.tfvars
postgres_db_name = "taskapp"
postgres_db_port = 5432
postgres_db_user = "postgresUser"
postgres_db_password="postgresDefault"
application_port = 3000
webapp_port = 80
auth0_domain = "{your-domain}.auth0.com"
auth0_api_identifier = ""
auth0_api_secret = ""
auth0_webapp_client_id = ""
auth0_anon_client_id = ""
auth0_anon_client_secret = ""
enable_e2e_dev_login = true
auth0_e2e_dev_client_id = ""
auth0_e2e_dev_client_subject = "{auth0_e2e_dev_client_id}@clients"
auth0_e2e_dev_client_secret = ""