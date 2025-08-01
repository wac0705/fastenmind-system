# N8N Integration for FastenMind

This directory contains N8N workflow definitions and integration documentation for the FastenMind quotation system.

## Overview

N8N is integrated into FastenMind to provide workflow automation capabilities for:
- Automated notifications
- Scheduled tasks
- Data synchronization
- Process automation
- Integration with external services

## Setup

### 1. Start N8N with Docker Compose

```bash
docker-compose up -d n8n
```

N8N will be available at: http://localhost:5678
- Username: admin
- Password: fastenmind123

### 2. Import Workflows

1. Access N8N at http://localhost:5678
2. Go to Workflows section
3. Import the JSON files from `./workflows/` directory

### 3. Configure Credentials

In N8N, configure the following credentials:

#### SMTP (Email)
- Type: SMTP
- Host: Your SMTP server
- Port: 587 (or your SMTP port)
- User: Your email
- Password: Your password
- SSL/TLS: Enable as needed

#### Slack (Optional)
- Type: Slack
- Access Token: Your Slack OAuth token
- Or use Webhook URL for simpler integration

#### API Authentication
- Type: Header Auth
- Header Name: Authorization
- Header Value: Bearer YOUR_API_TOKEN

## Available Workflows

### 1. Inquiry Notification (`inquiry-notification.json`)
**Trigger**: Webhook - When new inquiry is created
**Actions**:
- Fetch inquiry details from API
- Send email notification to assigned engineer
- Post notification to Slack channel
- Log the notification event

### 2. Daily Exchange Rate Update (`daily-exchange-rate.json`)
**Trigger**: Cron - Daily at 8 AM (weekdays)
**Actions**:
- Fetch latest exchange rates from external API
- Update rates in database
- Notify finance team via Slack
- Send alert email if significant changes detected

### 3. Quote Approval Workflow
**Trigger**: Event - When quote submitted for review
**Actions**:
- Check quote value against approval limits
- Route to appropriate approver
- Send approval request email
- Update quote status after approval

### 4. Cost Update Reminder
**Trigger**: Cron - Monthly on 1st at 9 AM
**Actions**:
- Check last update date for material costs
- Check last update date for process costs
- Send reminder emails to responsible engineers
- Create tasks in system for updates

### 5. Customer Credit Check
**Trigger**: Event - Before quote creation
**Actions**:
- Check customer credit limit
- Calculate outstanding amount
- Alert if over credit limit
- Require manager approval if needed

## Integration Points

### From FastenMind to N8N

The system triggers N8N workflows through:

1. **Direct Webhook Calls**
```javascript
await n8nService.triggerWorkflow({
  workflow_id: 'inquiry_notification',
  data: {
    event: 'inquiry.created',
    inquiry_id: inquiryId,
    timestamp: new Date().toISOString()
  }
})
```

2. **Event-Based Triggers**
Events logged in the system automatically trigger configured workflows:
- inquiry.created
- quote.created
- quote.submitted_for_review
- quote.approved
- customer.credit_limit_exceeded

3. **Scheduled Tasks**
Configured in the system and executed by N8N:
- Daily exchange rate updates
- Weekly report generation
- Monthly cost update reminders

### From N8N to FastenMind

N8N can interact with FastenMind through:

1. **REST API Calls**
- GET /api/inquiries/{id}
- POST /api/quotes/{id}/approve
- PUT /api/exchange-rates/batch-update

2. **Webhook Endpoints**
- POST /api/n8n/webhook/{event_type}

## Webhook URLs

After importing workflows, note the webhook URLs:

- Inquiry Notification: `http://localhost:5678/webhook/inquiry_notification`
- Quote Approval: `http://localhost:5678/webhook/quote_approval`
- Credit Check: `http://localhost:5678/webhook/credit_check`

## Environment Variables

Configure these in N8N:

```
API_URL=http://backend:8080
API_TOKEN=your-api-token
APP_URL=http://localhost:3000
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

## Monitoring

### Execution History
View all workflow executions in N8N:
1. Go to Executions tab
2. Filter by workflow, status, or date
3. Click on execution to see details

### From FastenMind
View N8N integration status:
1. Go to Workflows page
2. Check execution history
3. View any error messages

## Troubleshooting

### Common Issues

1. **Webhook not triggering**
   - Check webhook URL is correct
   - Verify N8N is running and accessible
   - Check network connectivity between services

2. **API authentication failing**
   - Verify API token is valid
   - Check token has necessary permissions
   - Ensure Authorization header format is correct

3. **Email not sending**
   - Verify SMTP credentials
   - Check firewall/security settings
   - Test SMTP connection in N8N

### Debug Mode

Enable debug mode in N8N:
1. Go to Settings
2. Enable "Save Execution Progress"
3. Enable "Save Data Progress"
4. Check execution details for each node

## Best Practices

1. **Error Handling**
   - Always include error handling nodes
   - Send notifications on workflow failures
   - Log errors for debugging

2. **Performance**
   - Use batch operations where possible
   - Implement rate limiting for external APIs
   - Set appropriate timeouts

3. **Security**
   - Use environment variables for sensitive data
   - Implement proper authentication
   - Validate webhook payloads

4. **Maintenance**
   - Regularly review execution logs
   - Update workflows as business logic changes
   - Document any custom nodes or functions

## Custom Workflows

To create custom workflows:

1. Design workflow in N8N interface
2. Test thoroughly with sample data
3. Export as JSON
4. Add to `./workflows/` directory
5. Document in this README

## Support

For issues or questions:
1. Check N8N documentation: https://docs.n8n.io
2. Review FastenMind API documentation
3. Contact system administrator