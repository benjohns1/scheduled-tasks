variable "env" {
    type = string
    default = "prod"
}

variable "aws_region" {
    type = string
    default = "us-west-2"
}

variable "aws_profile" {
    type = string
    default = ""
}

variable "application_name" {
    type = string
    default = "st"
}

variable "postgres_db_name" {
    type = string
}

variable "postgres_db_port" {
    type = number
}

variable "postgres_db_user" {
    type = string
}

variable "postgres_db_password" {
    type = string
}

variable "application_port" {
    type = number
}

variable "auth0_domain" {
    type = string
}

variable "auth0_api_identifier" {
    type = string
}

variable "auth0_api_secret" {
    type = string
}

variable "dbconn_maxretryattempts" {
    default = 20
    type = number
}

variable "dbconn_retrysleepseconds" {
    default = 3
    type = number
}