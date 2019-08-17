variable "aws_route53_zone" {
  type = string
}

variable "aws_route53_subdomain" {
  type = string
}

data "aws_route53_zone" "primary" {
    name = var.aws_route53_zone
}

resource "aws_route53_record" "subdomain" {
    zone_id = data.aws_route53_zone.primary.zone_id
    name = "${var.aws_route53_subdomain}.${data.aws_route53_zone.primary.name}"
    type = "A"
    ttl = "300"
    records = [aws_eip.ec2_eip.public_ip]
}
