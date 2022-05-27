# If you add an import that doesn't come with python 3.9, add it to requirements.txt
# If you make any changes to this file, be sure to re-zip the health files.  

# You can run `make health` to do this.

# DO NOT ZIP THE HEALTH FOLDER, zip the contents of the health folder, so there is no nested folder in the zip file, 
# and place it in the root of the api directory as health.zip

import requests
import os
import logging
import time

def check_health(req):    
  base_url = os.environ.get('VAMA_API_BASE_URL', 'err')
  
  if base_url == 'err':
    logging.fatal("No environment var for VAMA_API_BASE_URL")

  full_url = base_url + '/monitoring/health/db'
  timeoutSeconds = 5

  try:
    vama_api_response = requests.get(url = full_url, timeout = timeoutSeconds)
  except Exception as e:
    text = 'Exception occurred when making GET to {0}. Exception: {1}'.format(full_url, e)
    send_telegram_bot_message(text)
    return str(500)

  if vama_api_response.status_code != 200:
    time.sleep(10)
    try:
      retry_vama_api_response = requests.get(url = full_url, timeout = timeoutSeconds)
    except Exception as e:
      text = 'Exception occurred when making GET to {0}. Exception: {1}'.format(full_url, e)
      send_telegram_bot_message(text)
      return str(500)

    if retry_vama_api_response == 200:
      text = 'Health ping failed with code {0} but retry succeeded. Original failure: {1}'.format(vama_api_response.status_code, vama_api_response.text)
      return str(retry_vama_api_response.status_code)
    text = '{0} returned a {1} error. Message: {2}'.format(full_url, vama_api_response.status_code, vama_api_response.text)
    send_telegram_bot_message(text)
    
  return str(vama_api_response.status_code)

def send_telegram_bot_message(text): 
  telegram_api_token = os.environ.get('TELEGRAM_HEALTH_CHECKER_API_TOKEN', '')
  requests.post(
    url = 'https://api.telegram.org/bot{0}/sendMessage'.format(telegram_api_token),
    data = {
      'chat_id': -647388574, 
      'text': text
    }
  )