import boto3
import json


sqs = boto3.resource('sqs',
                        endpoint_url='http://localhost:9324',
                        region_name='elasticmq',
                        aws_secret_access_key='x',
                        aws_access_key_id='x',
                        use_ssl=False)


# Send Message to a queue
# Get the queue by its name
queue = sqs.get_queue_by_name(QueueName='ledger-salaried-earnings-job')

# Define your structured data as a Python dictionary
# Replace 'your_provider_salary_setting' with the actual data
provider_salary_setting = {
    'ProviderId': 123,  # example provider ID
    'DelaySalariedEarningsDayCount': 1,
}

sqs_body_dict = {
    'ProviderSalarySetting': provider_salary_setting,
    'LogLevel': "info"  # or the appropriate log level if applicable
}

sqs_body_json = json.dumps(sqs_body_dict)

# Send a message to the queue
response = queue.send_message(MessageBody=sqs_body_json)

print("Message sent! Message ID:", response.get('MessageId'))