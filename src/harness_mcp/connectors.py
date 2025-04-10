import requests
import argparse
import json
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

def list_connectors(connector_names=None, connector_identifiers=None, types=None):
    """
    List connectors from Harness API
    
    Args:
        connector_names: Optional list of connector names
        connector_identifiers: Optional list of connector identifiers
        types: Optional list of connector types
    
    Returns:
        API response as JSON
    """
    url = "https://app.harness.io/ng/api/connectors/listV2"
    
    # Query parameters
    params = {
        "accountIdentifier": "wFHXHD0RRQWoO8tIZT5YVw",
        "orgIdentifier": "default",
        "projectIdentifier": "nitisha",
        "getDefaultFromOtherRepo": "true",
        "getDistinctFromBranches": "true",
        "onlyFavorites": "false",
        "pageIndex": "0",
        "pageSize": "50",
        "sortOrders": "orderType=ASC"
    }
    
    # Headers
    headers = {
        "Content-Type": "application/json",
        "x-api-key": os.getenv("HARNESS_API_KEY")
    }
    
    # Request body
    payload = {
        "categories": ["CLOUD_PROVIDER"],
        "filterType": "Connector"
    }
    
    # Add optional parameters if provided
    if connector_names:
        payload["connectorNames"] = connector_names
    
    if connector_identifiers:
        payload["connectorIdentifiers"] = connector_identifiers
    
    if types:
        payload["types"] = types
    
    # Make the request
    response = requests.post(url, params=params, headers=headers, json=payload)
    
    return response.json()

def main():
    parser = argparse.ArgumentParser(description="List connectors from Harness API")
    parser.add_argument("--names", nargs="+", help="List of connector names")
    parser.add_argument("--identifiers", nargs="+", help="List of connector identifiers")
    parser.add_argument("--types", nargs="+", help="List of connector types")
    args = parser.parse_args()
    
    response = list_connectors(
        connector_names=args.names,
        connector_identifiers=args.identifiers,
        types=args.types
    )
    
    print(json.dumps(response, indent=2))

if __name__ == "__main__":
    main()
