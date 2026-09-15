resource "aws_elasticache_subnet_group" "this" {
  name       = "${var.name_prefix}-redis-subnet-group"
  subnet_ids = var.subnet_ids
  tags       = var.tags
}

resource "aws_security_group" "this" {
  name        = "${var.name_prefix}-redis-sg"
  description = "Acesso ao ElastiCache Redis somente a partir da VPC/EKS"
  vpc_id      = var.vpc_id

  ingress {
    description     = "Redis a partir dos security groups permitidos"
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = var.allowed_security_group_ids
    cidr_blocks     = length(var.allowed_security_group_ids) == 0 ? [var.vpc_cidr_block] : null
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, { Name = "${var.name_prefix}-redis-sg" })
}

resource "aws_elasticache_replication_group" "this" {
  replication_group_id       = "${var.name_prefix}-redis"
  description                = "Redis cache para targeting/evaluation do ToggleMaster"
  engine                     = "redis"
  node_type                  = var.node_type
  num_cache_clusters         = 1
  automatic_failover_enabled = false
  subnet_group_name          = aws_elasticache_subnet_group.this.name
  security_group_ids         = [aws_security_group.this.id]
  at_rest_encryption_enabled = true
  transit_encryption_enabled = true

  tags = var.tags
}
