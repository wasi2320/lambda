import botocore.exceptions
from moviepy.editor import VideoFileClip
import logging
import os
import json
import boto3

logger = logging.getLogger()
logger.setLevel("INFO")
s3 = boto3.client('s3', region_name='us-east-1')
dynamodb = boto3.resource('dynamodb', region_name='us-east-1')
dynamodb_table_name = 'processing_tracker'


def lambda_handler(event, context):
    bucket = event['bucket']
    key = event['key']
    table = dynamodb.Table(dynamodb_table_name)
    try:
        put_item_in_dynamodb(table, key, processed=0, error=None)
        resp = compares_video(bucket, key)
        put_item_in_dynamodb(table, key, processed=1, error=None if resp['statusCode'] == 200 else resp['body'])
        return resp
    except Exception as e:
        logger.error("error executing the function", e)
        put_item_in_dynamodb(table, key, processed=0, error=str(e))
        return {
            'statusCode': 500,
            'message': json.dumps({"message": "cannot run the func", "error": e})
        }


def compares_video(bucket, key):
    logger.info("Retrieving video from S3")
    filename = key.split('/')[1]
    try:
        logger.info("filename %s", filename)
        os.chmod("/tmp", 0o700)
        os.chdir("/tmp/")
        logger.info(os.getcwd())
        logger.info("change directory")
        s3.download_file(Bucket=bucket, Key=key, Filename=filename)
        resized_clip = reduce_video_size(filename)
        resized_clip.write_videofile(filename)
        s3.upload_file(filename, bucket, key)
        logger.info("Video resized and uploaded back to S3")
        return {
            'statusCode': 200,
            'body': json.dumps({"message": "Video processed successfully"})
        }
    except Exception as e:
        logger.error("Error processing video: %s", e)
        return {
            'statusCode': 500,
            'body': json.dumps({"error": str(e)})
        }


def reduce_video_size(video_input):
    logger.info("reducing the size of video from S3")
    clip = VideoFileClip(video_input)
    width_of_video = clip.w
    height_of_video = clip.h
    logger.info(f'Width and Height of original video: {width_of_video}x{height_of_video}')
    clip_resized = clip.resize(0.6)
    width_of_video = clip_resized.w
    height_of_video = clip_resized.h
    logger.info(f'Width and Height of original video: {width_of_video}x{height_of_video}')
    return clip_resized


def put_item_in_dynamodb(table, key, processed, error):
    logger.info("Updating the DynamoDB Table (put_item)")
    try:
        if error:
            item = {
                'object_key': key,
                'processed': processed,
                'error': error
            }
        else:
            item = {
                'object_key': key,
                'processed': processed,
            }

        table.put_item(Item=item)
        logger.info("Dynamo Table Updated")
    except botocore.exceptions.ClientError as e:
        logger.error("Error updating DynamoDB table: %s", e)
