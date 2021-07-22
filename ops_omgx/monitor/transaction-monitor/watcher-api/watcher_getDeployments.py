import json
import yaml
import pymysql
import boto3
import string
import random
import time
import requests
import redis

def watcher_getTransaction(event, context):

  # Parse incoming event
  body = json.loads(event["body"])
  address = body.get("address")
  
  # Read YML
  with open("env.yml", 'r') as ymlfile:
    config = yaml.load(ymlfile)

  # Get MySQL host and port
  endpoint = config.get('RDS_ENDPOINT')
  user = config.get('RDS_MYSQL_NAME')
  dbpassword = config.get('RDS_MYSQL_PASSWORD')
  dbname = config.get('RDS_DBNAME')

  con = pymysql.connect(endpoint, user=user, db=dbname,
                        passwd=dbpassword, connect_timeout=5)

  transactionData = []
  with con:
    try:
      cur = con.cursor()
      cur.execute("""SELECT hash, blockNumber, `from`, `to`, timestamp, crossDomainMessage FROM receipt WHERE `from`=%s ORDER BY CAST(blockNumber as unsigned) DESC""", (address))
      transactionsDataRaw = cur.fetchall()
      for transactionDataRaw in transactionsDataRaw:
        transactionData.append({
          "hash": transactionDataRaw[0],
          "blockNumber": int(transactionDataRaw[1]),
          "from": transactionDataRaw[2],
          "to": transactionDataRaw[3],
          "timeStamp": transactionDataRaw[4],
          "crossDomainMessage": transactionDataRaw[5]
        })
    except:
      transactionData = []

  con.close()
  
  response = {
    "statusCode": 201,
    "headers": {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Credentials": True,
      "Strict-Transport-Security": "max-age=63072000; includeSubdomains; preload",
      "X-Content-Type-Options": "nosniff",
      "X-Frame-Options": "DENY",
      "X-XSS-Protection": "1; mode=block",
      "Referrer-Policy": "same-origin",
      "Permissions-Policy": "*",
    },
    "body": json.dumps(transactionData),
  }
  return response

res = watcher_getTransaction({"body": json.dumps({
  "address": "0xb2780bABBe5Eaf6b611cAcC5cf3Db1C669224F60",
  "fromRange": "1",
  "toRange": "2"
})}, "context")

print(res)