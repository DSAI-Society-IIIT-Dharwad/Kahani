#!/usr/bin/env python3
"""
Simple test to verify Milvus Lite connection
"""

from pymilvus import MilvusClient
import os
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def test_milvus_lite():
    """Test different Milvus Lite connection methods"""
    
    # Test different URI formats
    uri_formats = [
        "./test_milvus.db",
        "test_milvus.db", 
        os.path.abspath("./test_milvus.db")
    ]
    
    for uri in uri_formats:
        try:
            logger.info(f"Testing URI: {uri}")
            client = MilvusClient(uri=uri)
            
            # Test basic operations
            collection_name = "test_collection"
            
            # Try to create a simple collection
            if client.has_collection(collection_name):
                client.drop_collection(collection_name)
            
            client.create_collection(
                collection_name=collection_name,
                dimension=384,
                metric_type="L2"
            )
            
            logger.info(f"✅ Success with URI: {uri}")
            
            # Cleanup
            client.drop_collection(collection_name)
            client.close()
            
            return uri
            
        except Exception as e:
            logger.error(f"❌ Failed with URI {uri}: {e}")
            continue
    
    return None

if __name__ == "__main__":
    result = test_milvus_lite()
    if result:
        print(f"✅ Milvus Lite works with URI: {result}")
    else:
        print("❌ Milvus Lite failed with all URI formats")
