# 📖 Kahani Documentation Index

**Complete documentation for the Kahani collaborative storytelling platform**

---

## 🗺️ Quick Navigation

| Document | Purpose | Audience |
|----------|---------|----------|
| **[README.md](../README.md)** | Project overview, quick start, tech stack | Everyone |
| **[ARCHITECTURE.md](ARCHITECTURE.md)** | System design, component breakdown, deployment | Developers, Architects |
| **[DATA_FLOW.md](DATA_FLOW.md)** | DFD diagrams, sequence diagrams, data models | Developers, System Analysts |
| **[BLOCKCHAIN.md](BLOCKCHAIN.md)** | Go blockchain implementation, PBFT, NFTs | Blockchain Developers |
| **[DIAGRAMS.md](DIAGRAMS.md)** | Mermaid.js visual diagrams | Visual Learners |
| **[Backend README](../Kahani_Ai_backend/README.md)** | FastAPI setup, API reference | Backend Developers |

---

## 📚 Documentation by Role

### For **New Contributors**
1. Start with [README.md](../README.md) - Project overview
2. Read [ARCHITECTURE.md](ARCHITECTURE.md) - System design
3. View [DIAGRAMS.md](DIAGRAMS.md) - Visual understanding
4. Check component READMEs for setup

### For **Frontend Developers**
1. [README.md](../README.md#quick-start) - Setup Next.js app
2. [ARCHITECTURE.md](ARCHITECTURE.md#frontend-layer) - Frontend architecture
3. [DATA_FLOW.md](DATA_FLOW.md#story-creation-sequence-diagram) - API flows
4. [DIAGRAMS.md](DIAGRAMS.md#story-creation-sequence-diagram) - UI sequences

### For **Backend Developers**
1. [Backend README](../Kahani_Ai_backend/README.md) - FastAPI setup
2. [ARCHITECTURE.md](ARCHITECTURE.md#ai-backend-layer) - RAG pipeline details
3. [DATA_FLOW.md](DATA_FLOW.md#dfd-level-2-story-suggestion-process-p1) - Data flows
4. [DIAGRAMS.md](DIAGRAMS.md#lore-extraction-flow) - Background tasks

### For **Blockchain Engineers**
1. [BLOCKCHAIN.md](BLOCKCHAIN.md) - Complete Go implementation
2. [ARCHITECTURE.md](ARCHITECTURE.md#blockchain-layer-work-in-progress) - Integration points
3. [DATA_FLOW.md](DATA_FLOW.md#dfd-level-2-nft-minting-p7---future-implementation) - NFT flows
4. [DIAGRAMS.md](DIAGRAMS.md#pbft-consensus-flow) - PBFT sequences

### For **DevOps/SRE**
1. [ARCHITECTURE.md](ARCHITECTURE.md#deployment) - Deployment configs
2. [BLOCKCHAIN.md](BLOCKCHAIN.md#deployment) - Multi-node setup
3. [README.md](../README.md#-technology-stack) - Infrastructure stack
4. [DIAGRAMS.md](DIAGRAMS.md#deployment-architecture) - Network topology

---

## 🔍 Documentation by Topic

### Architecture & Design
- [System Overview](ARCHITECTURE.md#system-overview)
- [Three-Tier Architecture](ARCHITECTURE.md#layer-breakdown)
- [Component Diagrams](DIAGRAMS.md#component-architecture)
- [Network Topology](DIAGRAMS.md#network-topology)

### Data Management
- [Database Schemas](ARCHITECTURE.md#database-architecture)
- [Data Flow Diagrams](DATA_FLOW.md)
- [ER Diagrams](DIAGRAMS.md#data-model-er-diagram)
- [Storage Architecture](ARCHITECTURE.md#storage-architecture)

### AI & Machine Learning
- [RAG Pipeline](../Kahani_Ai_backend/README.md#system-architecture)
- [Vector Search](ARCHITECTURE.md#vector-search-milvus)
- [LLM Integration](ARCHITECTURE.md#llm-api-groq-cloud)
- [Lore Extraction](DATA_FLOW.md#dfd-level-2-lore-extraction-p5)

### Blockchain & NFTs
- [PBFT Consensus](BLOCKCHAIN.md#consensus-mechanism)
- [Wallet System](BLOCKCHAIN.md#wallet-system)
- [NFT Implementation](BLOCKCHAIN.md#nft-implementation)
- [IPFS Integration](BLOCKCHAIN.md#ipfs-integration)

### API Reference
- [Frontend API Client](../lib/kahani-api.ts)
- [Backend Endpoints](../Kahani_Ai_backend/README.md#-api-endpoints)
- [Blockchain API](BLOCKCHAIN.md#api-reference)

### Security
- [Security Architecture](ARCHITECTURE.md#security-architecture)
- [Encryption Details](BLOCKCHAIN.md#wallet-system)
- [Authentication](ARCHITECTURE.md#authentication--authorization)

---

## 📊 Diagrams Quick Reference

### High-Level Diagrams
- [System Context](DIAGRAMS.md#system-context-diagram)
- [Component Architecture](DIAGRAMS.md#component-architecture)
- [Technology Stack](DIAGRAMS.md#technology-stack-layers)
- [Deployment](DIAGRAMS.md#deployment-architecture)

### Sequence Diagrams
- [Story Creation](DIAGRAMS.md#story-creation-sequence-diagram)
- [Lore Extraction](DIAGRAMS.md#lore-extraction-flow)
- [NFT Minting](DIAGRAMS.md#nft-minting-sequence-future)
- [PBFT Consensus](DIAGRAMS.md#pbft-consensus-flow)

### Data Flow Diagrams
- [Level 0 Context](DATA_FLOW.md#dfd-level-0-context-diagram)
- [Level 1 Processes](DATA_FLOW.md#dfd-level-1-main-system-processes)
- [Level 2 Details](DATA_FLOW.md#dfd-level-2-story-suggestion-process-p1)

### State Diagrams
- [Story Lifecycle](DIAGRAMS.md#state-transition-story-line-lifecycle)
- [PBFT States](DIAGRAMS.md#state-transition-pbft-block-states)

---

## 🚀 Getting Started Paths

### Path 1: Run Locally (AI-Only)
1. Clone repo: See [README.md](../README.md#1-clone-repository)
2. Start backend: See [Backend README](../Kahani_Ai_backend/README.md#-installation--setup)
3. Start frontend: See [README.md](../README.md#3-start-frontend)
4. Test system: See [Backend README](../Kahani_Ai_backend/README.md#-testing-the-system)

### Path 2: Deploy to Production
1. Review [Architecture](ARCHITECTURE.md#deployment)
2. Setup infrastructure: [Deployment Guide](ARCHITECTURE.md#production-deployment)
3. Configure networking: [Network Topology](DIAGRAMS.md#network-topology)
4. Monitor: [Observability](ARCHITECTURE.md#monitoring--observability)

### Path 3: Build Blockchain Layer
1. Read [Blockchain Specs](BLOCKCHAIN.md)
2. Follow [Implementation Phases](BLOCKCHAIN.md#implementation-phases)
3. Setup validators: [Deployment](BLOCKCHAIN.md#deployment)
4. Test consensus: [Testing Guide](BLOCKCHAIN.md#testing)

---

## 🛠️ Development Workflows

### Adding a New Feature
1. Review [Architecture](ARCHITECTURE.md) for integration points
2. Check [Data Flow](DATA_FLOW.md) for data transformations
3. Update [Diagrams](DIAGRAMS.md) if architecture changes
4. Document API changes in component READMEs

### Debugging Issues
1. Check [Troubleshooting](../Kahani_Ai_backend/README.md#-troubleshooting)
2. Review [Sequence Diagrams](DIAGRAMS.md#sequence-diagrams)
3. Trace data flow in [DATA_FLOW.md](DATA_FLOW.md)
4. Check logs as per [Monitoring](ARCHITECTURE.md#monitoring--observability)

### Performance Optimization
1. Review [Performance Metrics](ARCHITECTURE.md#performance-considerations)
2. Check [Scalability](ARCHITECTURE.md#scalability--horizontal-scaling)
3. Analyze [Data Access Patterns](DATA_FLOW.md#data-access-patterns)
4. Optimize based on [Bottlenecks](DATA_FLOW.md#performance-metrics)

---

## 📖 Documentation Standards

### Updating Documentation
When making changes to the system:

1. **Code Changes** → Update component README
   - Example: New API endpoint → Update [Backend README](../Kahani_Ai_backend/README.md)

2. **Architecture Changes** → Update [ARCHITECTURE.md](ARCHITECTURE.md)
   - Example: New microservice → Add to layer breakdown

3. **Data Model Changes** → Update [DATA_FLOW.md](DATA_FLOW.md)
   - Example: New database table → Add to data dictionary

4. **Visual Changes** → Update [DIAGRAMS.md](DIAGRAMS.md)
   - Example: New service → Add to component diagram

### Diagram Creation
1. Use [Mermaid.js](https://mermaid.js.org/) for all diagrams
2. Add diagrams to [DIAGRAMS.md](DIAGRAMS.md)
3. Reference diagrams in other docs using links
4. Test rendering in GitHub and VS Code

---

## 🔗 External Resources

### Technologies Used
- [Next.js Documentation](https://nextjs.org/docs)
- [FastAPI Documentation](https://fastapi.tiangolo.com/)
- [Milvus Documentation](https://milvus.io/docs)
- [Groq Cloud API](https://console.groq.com/)
- [IPFS Documentation](https://docs.ipfs.tech/)
- [Go Documentation](https://go.dev/doc/)

### Academic Papers
- [PBFT Paper](http://pmg.csail.mit.edu/papers/osdi99.pdf) - Byzantine Fault Tolerance
- [RAG Paper](https://arxiv.org/abs/2005.11401) - Retrieval-Augmented Generation
- [Ed25519 Spec](https://ed25519.cr.yp.to/) - Cryptographic signing

### Community
- GitHub Issues: Report bugs and request features
- Discord: Join discussions (coming soon)
- Twitter: Follow [@KahaniStories](https://twitter.com/kahanistories)

---

## 📝 Changelog

### Version 1.0 (Current - AI Phase)
- ✅ Next.js frontend with matte design
- ✅ FastAPI RAG backend
- ✅ Milvus vector search
- ✅ Story suggestion, editing, canonicalization
- ✅ Lore extraction
- ✅ PDF export
- ✅ Complete documentation suite

### Version 2.0 (Q2 2025 - Blockchain)
- 🚧 Go blockchain implementation
- 🚧 PBFT consensus (4 validators)
- 🚧 Ed25519 wallet auto-generation
- 🚧 BadgerDB storage
- 🚧 Supabase integration

### Version 3.0 (Q3 2025 - NFTs)
- 📋 IPFS integration
- 📋 NFT minting with co-authorship
- 📋 PageKite validator network
- 📋 Observer nodes
- 📋 Block explorer UI

---

## ❓ FAQ

**Q: Where do I start if I'm new to the project?**  
A: Read [README.md](../README.md) first, then [ARCHITECTURE.md](ARCHITECTURE.md), then component READMEs.

**Q: How do I run the system locally?**  
A: See [Backend README](../Kahani_Ai_backend/README.md#-installation--setup) and [Main README](../README.md#quick-start).

**Q: Where are the API docs?**  
A: Backend API: http://localhost:8000/docs (when running). Blockchain API: [BLOCKCHAIN.md](BLOCKCHAIN.md#api-reference).

**Q: How does PBFT consensus work?**  
A: See [BLOCKCHAIN.md](BLOCKCHAIN.md#consensus-mechanism) and [PBFT Flow Diagram](DIAGRAMS.md#pbft-consensus-flow).

**Q: Can I use this in production?**  
A: AI backend is production-ready. Blockchain layer is work in progress (Q2 2025 target).

**Q: How do I contribute?**  
A: Fork the repo, make changes, submit PR. See [Contributing](../README.md#-contributing).

---

## 📧 Support

- **Technical Questions**: Open a GitHub issue
- **Feature Requests**: Use issue templates
- **Bug Reports**: Include logs and reproduction steps
- **General Inquiries**: Contact via email (see [README.md](../README.md#-contact))

---

<div align="center">

**Kahani Documentation v1.0**

*Last Updated: January 2025*

[🏠 Home](../README.md) • [🏗️ Architecture](ARCHITECTURE.md) • [📊 Data Flow](DATA_FLOW.md) • [🔗 Blockchain](BLOCKCHAIN.md) • [📐 Diagrams](DIAGRAMS.md)

</div>
