# Job Sources Architecture

CareerOS separates discovery from provider implementation.

SourceConfig -> DiscoveryScheduler -> Adapter -> JobObservation -> Normalization/Dedup -> Canonical Job -> Analysis Queue

Adapters receive a source configuration and discovery query and return observations plus an optional provider cursor. They must not silently enrich source facts.

Provider implementations should prefer official/public APIs or explicit user-authorized integrations. Company career pages may be sources, but parsing must respect terms and technical access controls.

Each source has an independent cadence and rate limit. A discovery run is persisted before work begins and finalized with counts and status. Retries are bounded and idempotent.

Future providers can implement this contract without changing the canonical job model.