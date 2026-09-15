# IRSA: cada pod assume uma role com o mínimo de privilégio necessário,
# eliminando a necessidade de access keys estáticas dentro dos containers.

data "aws_iam_policy_document" "evaluation_assume_role" {
  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]
    principals {
      type        = "Federated"
      identifiers = [var.oidc_provider_arn]
    }
    condition {
      test     = "StringEquals"
      variable = "${var.oidc_provider_url}:sub"
      values   = ["system:serviceaccount:${var.namespace}:evaluation-sa"]
    }
  }
}

resource "aws_iam_role" "evaluation" {
  name               = "${var.name_prefix}-evaluation-irsa"
  assume_role_policy = data.aws_iam_policy_document.evaluation_assume_role.json
  tags               = var.tags
}

data "aws_iam_policy_document" "evaluation_permissions" {
  statement {
    sid       = "SendAnalyticsEvents"
    actions   = ["sqs:SendMessage", "sqs:GetQueueUrl"]
    resources = [var.sqs_queue_arn]
  }
}

resource "aws_iam_role_policy" "evaluation" {
  name   = "${var.name_prefix}-evaluation-policy"
  role   = aws_iam_role.evaluation.id
  policy = data.aws_iam_policy_document.evaluation_permissions.json
}

data "aws_iam_policy_document" "analytics_assume_role" {
  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]
    principals {
      type        = "Federated"
      identifiers = [var.oidc_provider_arn]
    }
    condition {
      test     = "StringEquals"
      variable = "${var.oidc_provider_url}:sub"
      values   = ["system:serviceaccount:${var.namespace}:analytics-sa"]
    }
  }
}

resource "aws_iam_role" "analytics" {
  name               = "${var.name_prefix}-analytics-irsa"
  assume_role_policy = data.aws_iam_policy_document.analytics_assume_role.json
  tags               = var.tags
}

data "aws_iam_policy_document" "analytics_permissions" {
  statement {
    sid       = "ConsumeAnalyticsEvents"
    actions   = ["sqs:ReceiveMessage", "sqs:DeleteMessage", "sqs:GetQueueUrl", "sqs:GetQueueAttributes"]
    resources = [var.sqs_queue_arn]
  }

  statement {
    sid       = "WriteAnalyticsTable"
    actions   = ["dynamodb:PutItem", "dynamodb:UpdateItem", "dynamodb:Query", "dynamodb:GetItem"]
    resources = [var.dynamodb_table_arn, "${var.dynamodb_table_arn}/index/*"]
  }
}

resource "aws_iam_role_policy" "analytics" {
  name   = "${var.name_prefix}-analytics-policy"
  role   = aws_iam_role.analytics.id
  policy = data.aws_iam_policy_document.analytics_permissions.json
}
