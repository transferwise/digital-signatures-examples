import os
import sys
import base64
import json
import rsa
import urllib3
from urllib.parse import urlencode
from datetime import datetime, timedelta, timezone

private_key_path = './private.pem'  # Change to private key path
profile_id = ''                     # Change to profile ID
account_id = ''                     # Change to borderless account ID
base_url = 'https://api.transferwise.com'

if os.getenv('API_TOKEN') is None:
    print('panic: no api token, please set with $ export API_TOKEN=xxx')
    sys.exit(0)
elif profile_id == '' or account_id == '':
    print('panic: profile / account ID missing, please add them')
    sys.exit(0)
elif os.path.exists(private_key_path) is False:
    print('panic: private key file not found, please update key path')
    sys.exit(0)

token = os.getenv('API_TOKEN')  # Set API token from env
http = urllib3.PoolManager()

def get_statement(one_time_token, signature):
    interval_start = (datetime.now(timezone.utc) - timedelta(days=14)).isoformat()
    interval_end = datetime.now(timezone.utc).isoformat()

    params = urlencode({
        'currency': 'GBP', 'type': 'COMPACT',
        'intervalStart': interval_start,
        'intervalEnd': interval_end})

    url = (
        base_url + '/v3/profiles/' + profile_id + '/borderless-accounts/' 
        + account_id + '/statement.json?' + params)

    headers = {
        'Authorization': 'Bearer ' + token,
        'User-Agent': 'tw-statements-sca',
        'Content-Type': 'application/json'}
    if one_time_token != "":
        headers['x-2fa-approval'] = one_time_token
        headers['X-Signature'] = signature
        print(headers['x-2fa-approval'], headers['X-Signature'])

    print('GET', url)

    r = http.request('GET', url, headers=headers, retries=False)

    print('status:', r.status)
    
    if r.status == 200 or r.status == 201:
        return json.loads(r.data)
    elif r.status == 403 and r.getheader('x-2fa-approval') is not None:
        one_time_token = r.getheader('x-2fa-approval')
        signature = do_sca_challenge(one_time_token)
        get_statement(one_time_token, signature)
    else:
        print('failed: ', r.status)
        print(r.data)
        sys.exit(0)

def do_sca_challenge(one_time_token):
    print('doing sca challenge')

    # Read the private key file as bytes.
    with open(private_key_path, 'rb') as f:
        private_key_data = f.read()

    private_key = rsa.PrivateKey.load_pkcs1(private_key_data, 'PEM')

    # Use the private key to sign the one-time-token that was returned 
    # in the x-2fa-approval header of the HTTP 403.
    signed_token = rsa.sign(
        one_time_token.encode('ascii'), 
        private_key, 
        'SHA-256')

    # Encode the signed message as friendly base64 format for HTTP 
    # headers.
    signature = base64.b64encode(signed_token).decode('ascii')

    return signature

def main():
    statement = get_statement("", "")

    if statement is not None and 'currency' in statement['request']:
        currency = statement['request']['currency']
    else:
        print('something is wrong')
        print(statement)
        sys.exit(0)
    
    if 'transactions' in statement:
        txns = len(statement['transactions'])
    else:
        print('Empty statement')
        sys.exit(0)

    print('\n', currency, 'statement received with', txns, 'transactions.')

if __name__ == '__main__':
    main()