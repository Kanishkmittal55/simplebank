import boto3
import time

# Connect to ElasticMQ
sqs = boto3.resource('sqs', endpoint_url='http://localhost:9324', region_name='us-east-1',
                     aws_access_key_id='dummy', aws_secret_access_key='dummy')

# List queues
for queue in sqs.queues.all():
    print("Queue URL:", queue.url)
    queue_name = queue.url.split('/')[-1]
    print("Queue Name:", queue_name)

    # Wait for 20 seconds (visibility timeout)
    time.sleep(20)

    # Receive messages from this specific queue
    queue = sqs.get_queue_by_name(QueueName=queue_name)
    messages = queue.receive_messages()
    if messages:
        print(f"Messages in {queue_name}:")
        for message in messages:
            print("Message Body:", message.body)
            # Uncomment below line to delete the message after reading
            # message.delete()
    else:
        print(f"No messages in {queue_name}")

