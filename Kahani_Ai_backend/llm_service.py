from groq import Groq
from config import get_settings
import logging
from typing import List, Optional
import json

logger = logging.getLogger(__name__)
settings = get_settings()


class LLMService:
    """Service for LLM operations using Groq"""
    
    def __init__(self):
        self.client = Groq(api_key=settings.groq_api_key)
        self.model = settings.llm_model
    
    def generate_story_suggestion(self, user_prompt: str, context: List[str]) -> str:
        """
        LLM-1: Generate story line suggestions based on user prompt and RAG context
        """
        context_text = "\n".join([f"- {c}" for c in context]) if context else "No previous context available."
        
        system_prompt = """You are a creative storytelling assistant. Based on the user's prompt and the story context provided, suggest a compelling next line or short passage (1-3 sentences max) that:
1. Flows naturally from the existing story
2. Addresses the user's request
3. Maintains consistency with established characters, settings, and plot
4. Is engaging and well-written

Return ONLY the suggested story line, nothing else."""
        
        user_message = f"""Story Context:
{context_text}

User Request: {user_prompt}

Suggest the next line(s):"""
        
        try:
            response = self.client.chat.completions.create(
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": user_message}
                ],
                model=self.model,
                temperature=0.7,
                max_tokens=200
            )
            
            suggestion = response.choices[0].message.content.strip()
            logger.info(f"Generated story suggestion: {suggestion[:50]}...")
            return suggestion
            
        except Exception as e:
            logger.error(f"Story suggestion generation failed: {e}")
            raise
    
    def extract_lore(self, story_lines: List[str]) -> dict:
        """
        Extract entities, characters, locations, and lore from story lines
        """
        story_text = "\n".join(story_lines)
        
        system_prompt = """You are a lore extraction specialist. Analyze the story text and extract:
1. Characters (with brief descriptions)
2. Locations (with brief descriptions)
3. Events (key plot points)
4. Items/Objects (significant items mentioned)

Return a JSON object with these categories."""
        
        user_message = f"""Extract lore from this story:

{story_text}

Return JSON format:
{{
  "characters": [{{"name": "...", "description": "..."}}, ...],
  "locations": [{{"name": "...", "description": "..."}}, ...],
  "events": [{{"name": "...", "description": "..."}}, ...],
  "items": [{{"name": "...", "description": "..."}}, ...]
}}"""
        
        try:
            response = self.client.chat.completions.create(
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": user_message}
                ],
                model=self.model,
                temperature=0.3,
                max_tokens=1000
            )
            
            lore_json = response.choices[0].message.content.strip()
            # Extract JSON from markdown code blocks if present
            if "```json" in lore_json:
                lore_json = lore_json.split("```json")[1].split("```")[0].strip()
            elif "```" in lore_json:
                lore_json = lore_json.split("```")[1].split("```")[0].strip()
            
            lore_data = json.loads(lore_json)
            logger.info(f"Extracted lore with {len(lore_data.get('characters', []))} characters")
            return lore_data
            
        except Exception as e:
            logger.error(f"Lore extraction failed: {e}")
            return {"characters": [], "locations": [], "events": [], "items": []}
    
    def summarize_story(self, story_lines: List[str]) -> str:
        """
        Generate a summary of story lines for embedding
        """
        story_text = "\n".join(story_lines)
        
        system_prompt = "You are a precise story summarizer. Create a concise summary that captures key plot points, characters, and themes. Keep it under 200 words."
        
        try:
            response = self.client.chat.completions.create(
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": f"Summarize this story:\n\n{story_text}"}
                ],
                model=self.model,
                temperature=0.3,
                max_tokens=300
            )
            
            summary = response.choices[0].message.content.strip()
            logger.info(f"Generated summary: {summary[:50]}...")
            return summary
            
        except Exception as e:
            logger.error(f"Summary generation failed: {e}")
            raise
    
    def canonicalize_story(self, story_lines: List[str]) -> str:
        """
        LLM-2: Create a canonical, polished version of the story
        """
        story_text = "\n".join([f"{i+1}. {line}" for i, line in enumerate(story_lines)])
        
        system_prompt = """You are a professional story editor. Transform the provided story lines into a polished, canonical narrative that:
1. Maintains all plot points and character actions
2. Improves flow and coherence
3. Fixes any inconsistencies
4. Enhances readability while preserving the author's voice
5. Formats it as a proper story (with paragraphs, not numbered lines)

Return the complete canonical story."""
        
        user_message = f"""Story lines to canonicalize:

{story_text}

Create the canonical version:"""
        
        try:
            response = self.client.chat.completions.create(
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": user_message}
                ],
                model=self.model,
                temperature=0.5,
                max_tokens=2000
            )
            
            canonical = response.choices[0].message.content.strip()
            logger.info(f"Generated canonical story ({len(canonical)} chars)")
            return canonical
            
        except Exception as e:
            logger.error(f"Canonicalization failed: {e}")
            raise


# Global instance
llm_service = LLMService()
