#!/usr/bin/env python3
"""
Example usage script for Kahani AI
Demonstrates the complete workflow
"""

import requests
import json
import time

BASE_URL = "http://localhost:8000"

def print_section(title):
    print("\n" + "="*60)
    print(f"  {title}")
    print("="*60)

def print_response(data, truncate=True):
    if isinstance(data, dict) and 'suggestion' in data and truncate:
        print(f"Suggestion: {data['suggestion'][:80]}...")
        print(f"Context items: {data.get('context_count', 0)}")
    else:
        print(json.dumps(data, indent=2)[:500])

def main():
    print("""
    ╔═══════════════════════════════════════════════════════════╗
    ║              Kahani AI - Example Workflow                 ║
    ║         RAG-Powered Story Generation System               ║
    ╚═══════════════════════════════════════════════════════════╝
    """)
    
    # 1. Health Check
    print_section("1. Health Check")
    try:
        response = requests.get(f"{BASE_URL}/health")
        print(f"✅ Status: {response.json()['status']}")
        print(f"✅ Milvus: {response.json()['milvus_connected']}")
        print(f"✅ Database: {response.json()['database_ok']}")
    except requests.exceptions.ConnectionError:
        print("❌ Cannot connect to server. Is it running?")
        print("\nStart the server with: python main.py")
        return
    
    # 2. First Story Line
    print_section("2. Creating First Story Line")
    print("Prompt: 'A wizard lives in a crystal tower'")
    
    response = requests.post(f"{BASE_URL}/api/story/suggest", json={
        "user_prompt": "A wizard lives in a crystal tower high above the clouds",
        "user_id": "demo_user"
    })
    suggestion1 = response.json()
    print_response(suggestion1)
    
    # Sign first line
    print("\n📝 Signing and storing...")
    response = requests.post(f"{BASE_URL}/api/story/edit", json={
        "llm_proposed": suggestion1['suggestion'],
        "final_text": "In a tower of gleaming crystal, high above the misty valleys, lived the ancient wizard Eldrin, keeper of forgotten secrets.",
        "user_id": "demo_user"
    })
    print(f"✅ Stored as line #{response.json()['id']}")
    time.sleep(2)  # Wait for embedding to be created
    
    # 3. Second Story Line (with context!)
    print_section("3. Creating Second Story Line (with RAG context)")
    print("Prompt: 'A young warrior seeks the wizard's help'")
    
    response = requests.post(f"{BASE_URL}/api/story/suggest", json={
        "user_prompt": "A young warrior arrives seeking the wizard's help against a dragon",
        "user_id": "demo_user"
    })
    suggestion2 = response.json()
    print_response(suggestion2)
    print(f"\n💡 Notice: Now has {suggestion2['context_count']} context items from previous story!")
    
    # Sign second line
    print("\n📝 Signing and storing...")
    response = requests.post(f"{BASE_URL}/api/story/edit", json={
        "llm_proposed": suggestion2['suggestion'],
        "final_text": "One stormy night, a young warrior named Kai climbed the winding stairs, desperate to seek Eldrin's counsel about the dragon terrorizing the kingdom.",
        "user_id": "demo_user"
    })
    print(f"✅ Stored as line #{response.json()['id']}")
    time.sleep(2)
    
    # 4. Add more lines
    print_section("4. Adding More Story Lines")
    
    lines = [
        "The wizard sensed dark magic stirring in the ancient mountain peaks.",
        "A thousand-year-old dragon named Zarthax had awakened from its slumber.",
        "The dragon sought revenge against those who sealed it away long ago."
    ]
    
    for i, line in enumerate(lines, start=3):
        print(f"\nLine {i}: {line[:50]}...")
        requests.post(f"{BASE_URL}/api/story/edit", json={
            "llm_proposed": line,
            "final_text": line,
            "user_id": "demo_user"
        })
        print("✅ Stored")
        time.sleep(1)
    
    # 5. View complete story
    print_section("5. Viewing Complete Story")
    response = requests.get(f"{BASE_URL}/api/story/lines?verified_only=true")
    lines = response.json()
    
    print(f"\n📖 Story so far ({len(lines)} lines):\n")
    for line in lines:
        print(f"{line['line_number']}. {line['text']}")
    
    # 6. Extract Lore
    print_section("6. Extracting Lore (Characters, Locations, Events)")
    line_ids = [line['id'] for line in lines]
    
    response = requests.post(f"{BASE_URL}/api/lore/extract", json={
        "line_ids": line_ids
    })
    lore = response.json()
    
    print(f"\n🧙 Characters ({len(lore['characters'])}):")
    for char in lore['characters']:
        print(f"  • {char['name']}: {char['description']}")
    
    print(f"\n🏰 Locations ({len(lore['locations'])}):")
    for loc in lore['locations']:
        print(f"  • {loc['name']}: {loc['description']}")
    
    print(f"\n⚔️ Events ({len(lore['events'])}):")
    for event in lore['events']:
        print(f"  • {event['name']}: {event['description']}")
    
    # 7. Test Context Retrieval
    print_section("7. Testing RAG Context Retrieval")
    print("Query: 'wizard tower'")
    
    response = requests.post(f"{BASE_URL}/api/context/retrieve", json={
        "query": "wizard tower crystal",
        "top_k": 3
    })
    results = response.json()
    
    print(f"\n🔍 Found {results['count']} relevant contexts:")
    for i, item in enumerate(results['results'][:3], 1):
        print(f"\n{i}. {item['text'][:80]}...")
        print(f"   Score: {item['score']:.3f} | Type: {item['content_type']}")
    
    # 8. Create Canonical Version
    print_section("8. Creating Canonical Story (LLM-2)")
    print("Transforming story lines into polished narrative...")
    
    response = requests.post(f"{BASE_URL}/api/story/canonicalize", json={
        "title": "The Dragon of Crystal Tower"
    })
    canonical = response.json()
    
    print(f"\n📚 {canonical['title']}")
    print(f"Original lines: {canonical['original_lines_count']}")
    print(f"\nCanonical Version:\n")
    print(canonical['full_text'][:400] + "...")
    print(f"\n✅ Canonical story ID: {canonical['id']}")
    
    # 9. Summary
    print_section("Summary")
    print("""
    ✅ Story Creation Flow Completed!
    
    What happened:
    1. Created 5 story lines
    2. Each line got better context from RAG
    3. Extracted characters, locations, and events
    4. Created canonical polished version
    5. All content embedded and searchable
    
    Background tasks (if running):
    • Every 30 min: Extract lore automatically
    • Every hour: Generate summaries
    
    Next steps:
    • Visit http://localhost:8000/ui for web interface
    • Check http://localhost:8000/docs for API docs
    • Continue building your story!
    """)

if __name__ == "__main__":
    main()
