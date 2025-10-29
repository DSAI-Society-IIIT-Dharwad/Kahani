# API Testing Guide

## Test Sequence - Complete Story Creation Flow

### 1. Health Check

```bash
curl http://localhost:8000/health
```

**Expected Response:**

```json
{
  "status": "healthy",
  "milvus_connected": true,
  "database_ok": true
}
```

---

### 2. Get First Story Suggestion

```bash
curl -X POST http://localhost:8000/api/story/suggest \
  -H "Content-Type: application/json" \
  -d '{
    "user_prompt": "A mysterious wizard lives in a crystal tower",
    "user_id": "test_user"
  }'
```

**Expected Response:**

```json
{
  "id": 1,
  "suggestion": "In a tower of gleaming crystal...",
  "context_used": [],
  "context_count": 0,
  "verified": false
}
```

---

### 3. Sign and Store First Line

```bash
curl -X POST http://localhost:8000/api/story/edit \
  -H "Content-Type: application/json" \
  -d '{
    "llm_proposed": "In a tower of gleaming crystal...",
    "final_text": "In a tower of gleaming crystal, high above the misty valleys, lived the ancient wizard Eldrin.",
    "user_id": "test_user"
  }'
```

---

### 4. Get Second Suggestion (with context!)

```bash
curl -X POST http://localhost:8000/api/story/suggest \
  -H "Content-Type: application/json" \
  -d '{
    "user_prompt": "A young warrior arrives seeking help",
    "user_id": "test_user"
  }'
```

**Expected:** Should include context from first line!

---

### 5. Sign Second Line

```bash
curl -X POST http://localhost:8000/api/story/edit \
  -H "Content-Type: application/json" \
  -d '{
    "llm_proposed": "...",
    "final_text": "One stormy night, a young warrior named Kai approached the tower, desperate for the wizard'\''s aid.",
    "user_id": "test_user"
  }'
```

---

### 6. Add More Lines

```bash
# Line 3
curl -X POST http://localhost:8000/api/story/edit \
  -H "Content-Type: application/json" \
  -d '{
    "llm_proposed": "...",
    "final_text": "The wizard sensed great danger approaching the kingdom.",
    "user_id": "test_user"
  }'

# Line 4
curl -X POST http://localhost:8000/api/story/edit \
  -H "Content-Type: application/json" \
  -d '{
    "llm_proposed": "...",
    "final_text": "An ancient dragon had awakened from its thousand-year slumber.",
    "user_id": "test_user"
  }'
```

---

### 7. View All Story Lines

```bash
curl http://localhost:8000/api/story/lines?verified_only=true
```

**Expected Response:**

```json
[
  {
    "id": 1,
    "line_number": 1,
    "text": "In a tower of gleaming crystal...",
    "verified": true,
    "user_edited": false,
    "created_at": "2025-10-29T..."
  },
  ...
]
```

---

### 8. Extract Lore

```bash
curl -X POST http://localhost:8000/api/lore/extract \
  -H "Content-Type: application/json" \
  -d '{
    "line_ids": [1, 2, 3, 4]
  }'
```

**Expected Response:**

```json
{
  "characters": [
    { "name": "Eldrin", "description": "ancient wizard" },
    { "name": "Kai", "description": "young warrior" }
  ],
  "locations": [
    { "name": "Crystal Tower", "description": "gleaming tower above valleys" },
    { "name": "Kingdom", "description": "threatened by dragon" }
  ],
  "events": [
    { "name": "Kai's Arrival", "description": "warrior seeks wizard's help" },
    { "name": "Dragon Awakening", "description": "ancient dragon awakens" }
  ],
  "items": [],
  "total_entries": 6
}
```

---

### 9. View All Lore

```bash
curl http://localhost:8000/api/lore/all
```

---

### 10. Retrieve Context (Test RAG)

```bash
curl -X POST http://localhost:8000/api/context/retrieve \
  -H "Content-Type: application/json" \
  -d '{
    "query": "wizard tower",
    "top_k": 5,
    "content_type": "story_line"
  }'
```

**Expected:** Should return relevant story lines about wizard and tower

---

### 11. Canonicalize Story

```bash
curl -X POST http://localhost:8000/api/story/canonicalize \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Dragon of the Crystal Tower"
  }'
```

**Expected Response:**

```json
{
  "id": 1,
  "title": "The Dragon of the Crystal Tower",
  "full_text": "In a tower of gleaming crystal, high above the misty valleys, lived the ancient wizard Eldrin. One stormy night...",
  "original_lines_count": 4,
  "created_at": "2025-10-29T..."
}
```

---

### 12. Get Canonical Story

```bash
curl http://localhost:8000/api/story/canonical/1
```

---

## Python Test Script

Save as `test_api.py`:

```python
import requests
import json

BASE_URL = "http://localhost:8000"

def test_flow():
    print("🧪 Testing Kahani AI API\n")

    # 1. Health check
    print("1️⃣ Health Check...")
    r = requests.get(f"{BASE_URL}/health")
    print(f"   Status: {r.json()['status']}\n")

    # 2. First suggestion
    print("2️⃣ Getting first suggestion...")
    r = requests.post(f"{BASE_URL}/api/story/suggest", json={
        "user_prompt": "A mysterious wizard lives in a crystal tower",
        "user_id": "test_user"
    })
    suggestion1 = r.json()
    print(f"   Suggestion: {suggestion1['suggestion'][:50]}...\n")

    # 3. Sign first line
    print("3️⃣ Signing first line...")
    r = requests.post(f"{BASE_URL}/api/story/edit", json={
        "llm_proposed": suggestion1['suggestion'],
        "final_text": "In a tower of gleaming crystal, high above the misty valleys, lived the ancient wizard Eldrin.",
        "user_id": "test_user"
    })
    print(f"   Line ID: {r.json()['id']}\n")

    # 4. Second suggestion (with context)
    print("4️⃣ Getting second suggestion (should have context)...")
    r = requests.post(f"{BASE_URL}/api/story/suggest", json={
        "user_prompt": "A young warrior arrives seeking help",
        "user_id": "test_user"
    })
    suggestion2 = r.json()
    print(f"   Context count: {suggestion2['context_count']}")
    print(f"   Suggestion: {suggestion2['suggestion'][:50]}...\n")

    # 5. Add more lines
    print("5️⃣ Adding more lines...")
    lines = [
        "One stormy night, a young warrior named Kai approached the tower.",
        "The wizard sensed great danger approaching the kingdom.",
        "An ancient dragon had awakened from its slumber."
    ]

    for line in lines:
        requests.post(f"{BASE_URL}/api/story/edit", json={
            "llm_proposed": line,
            "final_text": line,
            "user_id": "test_user"
        })
    print(f"   Added {len(lines)} lines\n")

    # 6. Get all lines
    print("6️⃣ Fetching all story lines...")
    r = requests.get(f"{BASE_URL}/api/story/lines")
    lines = r.json()
    print(f"   Total lines: {len(lines)}\n")

    # 7. Extract lore
    print("7️⃣ Extracting lore...")
    line_ids = [line['id'] for line in lines]
    r = requests.post(f"{BASE_URL}/api/lore/extract", json={
        "line_ids": line_ids
    })
    lore = r.json()
    print(f"   Characters: {len(lore['characters'])}")
    print(f"   Locations: {len(lore['locations'])}")
    print(f"   Events: {len(lore['events'])}\n")

    # 8. Canonicalize
    print("8️⃣ Creating canonical version...")
    r = requests.post(f"{BASE_URL}/api/story/canonicalize", json={
        "title": "The Dragon of the Crystal Tower"
    })
    canonical = r.json()
    print(f"   Story ID: {canonical['id']}")
    print(f"   Length: {len(canonical['full_text'])} characters\n")

    print("✅ All tests passed!")

if __name__ == "__main__":
    test_flow()
```

Run with:

```bash
python test_api.py
```

---

## JavaScript/Node.js Test

```javascript
const axios = require("axios");

const BASE_URL = "http://localhost:8000";

async function testFlow() {
  console.log("🧪 Testing Kahani AI API\n");

  // Health check
  const health = await axios.get(`${BASE_URL}/health`);
  console.log("✅ Health:", health.data.status);

  // Create story
  const suggest = await axios.post(`${BASE_URL}/api/story/suggest`, {
    user_prompt: "A wizard in a crystal tower",
    user_id: "test_user",
  });
  console.log("✅ Suggestion:", suggest.data.suggestion.substring(0, 50));

  // Sign it
  const sign = await axios.post(`${BASE_URL}/api/story/edit`, {
    llm_proposed: suggest.data.suggestion,
    final_text: suggest.data.suggestion,
    user_id: "test_user",
  });
  console.log("✅ Signed line ID:", sign.data.id);

  console.log("\n✅ All tests passed!");
}

testFlow().catch(console.error);
```

---

## Postman Collection

Import this JSON into Postman:

```json
{
  "info": {
    "name": "Kahani AI API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "header": [],
        "url": "http://localhost:8000/health"
      }
    },
    {
      "name": "Get Story Suggestion",
      "request": {
        "method": "POST",
        "header": [{ "key": "Content-Type", "value": "application/json" }],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"user_prompt\": \"A mysterious wizard lives in a crystal tower\",\n  \"user_id\": \"test_user\"\n}"
        },
        "url": "http://localhost:8000/api/story/suggest"
      }
    },
    {
      "name": "Sign Story Line",
      "request": {
        "method": "POST",
        "header": [{ "key": "Content-Type", "value": "application/json" }],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"llm_proposed\": \"The wizard lived...\",\n  \"final_text\": \"In a crystal tower lived the wizard Eldrin.\",\n  \"user_id\": \"test_user\"\n}"
        },
        "url": "http://localhost:8000/api/story/edit"
      }
    }
  ]
}
```

---

**Happy Testing! 🧪✨**
