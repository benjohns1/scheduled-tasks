resource "aws_alb_target_group" "webapp_alb_group" {
    name = "${var.prefix}-alb-group"
    port = var.host_webapp_port
    protocol = "HTTP"
    vpc_id = data.aws_vpc.default.id
}

resource "aws_alb" "alb_main" {
    name = "${var.prefix}-alb-ecs"
    subnets = data.aws_subnet_ids.all.ids
    security_groups = [aws_security_group.nsg_task.id]
}

resource "aws_alb_listener" "webapp_alb_listener" {
    load_balancer_arn = aws_alb.alb_main.arn
    port = var.host_webapp_port
    protocol = "HTTP"
    default_action {
        target_group_arn = aws_alb_target_group.webapp_alb_group.id
        type = "forward"
    }
}