# RemedyMate Backend API Specification (V1)

**Version:** 1.0
**Date:** Aug 25, 2025
**Contact:** Backend Team Lead

## 1. Overview

This document provides the official specification for the RemedyMate V1 backend API. The API is designed to be a simple, stateless, and secure service that provides content and AI-driven classification for the mobile and web clients.

All requests and responses will be in **JSON** format.

### General Information

*   **Base URL (Development):** `localhost:8080/api/v1`
*   **Authentication:** V1 of the API is unauthenticated and does not require an API key from the client.

## 2. Architectural Notes: Chat Sessions

A critical architectural decision for V1 is that the **backend is completely stateless**.

*   **No Server-Side Sessions:** The backend does not store any chat history or maintain user sessions. Each API call is an independent, atomic transaction.
*   **Client-Side Responsibility:** The **mobile client is responsible for managing the user's chat history**. As per the PRD, the last 10 chats should be saved locally on the user's device for privacy and offline access.

This approach simplifies the backend, enhances user privacy, and aligns with the project's core principles.

## 3. API Endpoints

Here are the details for each of the five core endpoints.

---

### 3.1 Triage Symptom Input

This is the first and most critical endpoint. It performs a safety check on the user's input to identify any red-flag conditions.

*   **Endpoint:** `POST /triage`
*   **Description:** Takes raw user text and returns a safety classification level (`GREEN`, `YELLOW`, or `RED`).

#### Request Body

```json
{
  "user_input": "I have a bad headache and feel tired.",
  "language": "en"
}
```

#### Response Body (200 OK)

**Example 1: Safe Input (GREEN)**
```json
{
  "level": "GREEN",
  "red_flags": []
}
```

**Example 2: Dangerous Input (RED)**
```json
{
  "level": "RED",
  "red_flags": [
    "chest pain"
  ]
}
```

---

### 3.2 Map Input to Topic

If the triage result is not `RED`, this endpoint maps the user's symptoms to a predefined content topic.

*   **Endpoint:** `POST /map-topic`
*   **Description:** Takes user text and returns the single most relevant `topic_key`.

#### Request Body

```json
{
  "user_input": "My stomach feels bloated after I eat.",
  "language": "en"
}
```

#### Response Body (200 OK)

```json
{
  "topic_key": "simple_indigestion"
}
```

---

### 3.3 Compose Guidance Card

This endpoint takes a `topic_key` and generates the final, user-facing guidance card in markdown format.

*   **Endpoint:** `POST /compose`
*   **Description:** Assembles a markdown-formatted advice card from the approved content blocks.

#### Request Body

```json
{
  "topic_key": "common_cold_adult",
  "language": "en"
}
```

#### Response Body (200 OK)

```json
{
  "markdown_card": "### Self-Care Ideas\n*   **Stay Hydrated:** Drink plenty of water, warm tea, or clear broths.\n*   **Get Plenty of Rest:** Your body needs energy to fight the infection.\n*   **Use a Saline Nasal Spray:** This can help relieve nasal congestion.\n\n### Over-the-Counter (OTC) Suggestions\n*   You may ask a pharmacist about a **decongestant** to help with a stuffy nose. Always read the label and follow their advice.\n\n### When to Seek Care\n*   If you have a high fever (above 38.5°C or 101.3°F).\n*   If your symptoms last for more than 10 days.\n*   If you have difficulty breathing.\n\n---\n**Disclaimer:** This is general information, not medical advice. If your symptoms are severe or persist, please see a healthcare professional."
}
```

---

### 3.4 Get Content for a Topic

Provides the raw, structured content for a specific topic, useful for caching or direct display.

*   **Endpoint:** `GET /content/:lang/:topic_key`
*   **Example URL:** `https://localhost:8080/api/v1/content/am/common_cold_adult`
*   **Description:** Fetches the approved blocks for a single topic in all available languages.

#### Response Body (200 OK)

```json
{
    "topic_key": "headache",
    "translations": {
      "am": {
        "self_care": [
          "በጸጥታ እና ጨለማ ክፍል ውስጥ ዕረፍት ይውሰዱ።",
          "ቀዝቃዛ እና እርጥብ ጨርቅ ግንባርዎ ላይ ይደግፉ።",
          "በቂ ውሃ ይጠጡ። በየጊዜው ካፌይን የሚጠጡ ከሆነ አነስተኛ መጠን ያለው ካፌይን (እንደ አንድ ኩባያ ቡና ወይም ሻይ) ሊረዳ ይችላል።"
        ],
        "otc_categories": [
          {
            "category_name": "ህመም ማስታገሻ",
            "safety_note": "ትክክለኛውን እንዲመርጡልዎ ፋርማሲስት ስለ ፓራሲታሞል ወይም አይቡፕሮፌን መጠየቅ ይችላሉ። ሁልጊዜ በጥቅሉ ላይ ያለውን የመድሀኒት መጠን ይከተሉ እና ከሚመከረው መጠን በላይ አይወስዱ።"
          }
        ],
        "seek_care_if": [
          "የራስ ህመምዎ በጣም ድንገተኛ እና ከባድ ነው (ብዙውን ጊዜ 'በሕይወትዎ ውስጥ እጅግ የከፋ የራስ ምታት' ወይም 'በጣም ከባድ የራስ ምታት' ተብሎ ይገለጻል)።",
          "ከፍተኛ ትኩሳት፣ ጠባሳ አንገት፣ ግርግር ወይም የእይታ ለውጥ አለዎት።",
          "ከቅርብ ጊዜ የጭንቅላት ጉዳት ተከስቷል።",
          "ህመሙ ከባድ ነው እና በእረፍት ወይም በገበያ ውስጥ በሚገኝ መድሀኒት የማይሻሻል ከሆነ።",
          "በሰውነትዎ በአንደኛው ጎን የመድከም ወይም የመደንዘዝ ስሜት፣ ወይም ለመናገር መቸገር ካጋጠመዎት።"
        ],
        "disclaimer": "ይህ አጠቃላይ መረጃ ነው፣ የሕክምና ምክር አይደለም። ምልክቶችዎ ከባድ ወይም ቆይተው ከሚታወቁ ከሆነ፣ እባክዎን የሕክምና ባለሙያ ይጠይቁ።"
      }
    }
  }
```

---

### 3.5 Get Offline Pack

This endpoint provides a single, bundled JSON object containing all the approved content for the Top 30 topics. The mobile client should call this once on first launch (with a good connection) and store the response locally.

*   **Endpoint:** `GET /offline-pack`
*   **Description:** Fetches all approved content for device caching.

#### Response Body (200 OK)

*(Note: This is a truncated example for brevity. The full response will contain all 30 topics.)*
```json
[
 {
   "topic_key": "indigestion",
   "translations": {
     "en": {
       "self_care": [
         "Eat smaller, more frequent meals throughout the day.",
         "Avoid foods that trigger your indigestion, such as fatty, spicy, or acidic foods.",
         "Drink peppermint or ginger tea to help settle your stomach.",
         "Loosen tight clothing around your waist.",
         "Avoid eating close to bedtime."
       ],
       "otc_categories": [
         {
           "category_name": "antacid",
           "safety_note": "You can ask a pharmacist about antacids to reduce stomach acid. Always follow the dosage instructions on the label and do not use them for too long without medical advice."
         }
       ],
       "seek_care_if": [
         "You experience severe or persistent stomach pain.",
         "You have black, tarry, or bloody stools.",
         "You have unexplained weight loss.",
         "You experience difficulty or pain when swallowing.",
         "Your indigestion is new, worsening, or accompanied by dizziness or sweating."
       ],
       "disclaimer": "This is general information, not medical advice. If your symptoms are severe or persistent, please consult a healthcare professional."
     },
     "am": {
       "self_care": [
         "በቀን ውስጥ ትንሽ እና ደጋግመው ይመገቡ።",
         "የሆድ ቁርጠትዎን የሚያስነሱ ምግቦችን ያስወግዱ፣ ለምሳሌ ቅባት የበዛባቸው፣ ቅመም የበዛባቸው ወይም አሲዳማ የሆኑ ምግቦች።",
         "ሆድዎን ለማረጋጋት ሚንት ወይም ዝንጅብል ሻይ ይጠጡ።",
         "በወገብዎ ላይ ያለውን ጥብቅ ልብስ ያላቅቁ።",
         "ከመኝታ ሰዓትዎ በፊት ምግብ ከመመገብ ይቆጠቡ።"
       ],
       "otc_categories": [
         {
           "category_name": "አንታሲድ (Antacid)",
           "safety_note": "የሆድ አሲድ ለመቀነስ ስለ አንታሲድ ፋርማሲስት መጠየቅ ይችላሉ። ሁልጊዜ በጥቅሉ ላይ ያለውን የመድሃኒት መጠን ይከተሉ እና የህክምና ምክር ሳይወስዱ ለረጅም ጊዜ አይጠቀሙባቸው።"
         }
       ],
       "seek_care_if": [
         "ከባድ ወይም የማያቋርጥ የሆድ ህመም ካለብዎት።",
         "ጥቁር፣ ሬንጅ የሚመስል ወይም ደም የያዘ ሰገራ ካለብዎት።",
         "ያልታወቀ ክብደት መቀነስ ካለብዎት።",
         "ምግብ ሲውጡ አስቸጋሪነት ወይም ህመም ካለብዎት።",
         "የሆድ ቁርጠትዎ አዲስ ከሆነ፣ እየባሰ ከሄደ፣ ወይም ከማዞር ወይም ከላብ ጋር አብሮ ከመጣ።"
       ],
       "disclaimer": "ይህ አጠቃላይ መረጃ ነው፣ የሕክምና ምክር አይደለም። ምልክቶችዎ ከባድ ወይም ቆይተው ከሚታወቁ ከሆነ፣ እባክዎን የሕክምና ባለሙያ ይጠይቁ።"
     }
   }
 },
 .
 .
 .
]
```