# STAR Stories

Behavioral interview prep, written in the STAR format
(Situation → Task → Action → Result).

## How to use this file

- **Memorize the skeleton, not the text.** In an interview, tell the
  story from the bullet points in your own words. A recited answer
  sounds robotic.
- **Target length: 1:45–2:10 spoken** (~260–300 words).
- **Use "I" for my actions.** Credit the team when it applies, but make
  clear which part was mine.
- **Every number must survive "how do you know?"** If I can't say where a
  number came from, it doesn't go in the story.
- **Details are follow-ups, not part of the story.** Keep the story at the
  level of decisions; each detail left out is a follow-up where I can
  go deep.

| # | Story | Main competencies | Status |
|---|---|---|---|
| 1 | Bank reconciliation (ERP) | Technical challenge, ownership, trade-offs | Done |
| 2 | Conflict or leadership (Apple Developer Academy) | Collaboration, leadership | To do |
| 3 | A failure and what I learned | Accountability, growth | To do |
| 4 | Tight deadline | Prioritization | To do |
| 5 | Disagreeing with someone on the team | Communication, influence | To do |
| 6 | Full ownership of a problem | Ownership | To do |

---

## Story #1: Bank reconciliation

**Answers questions like:**
- "What's the biggest technical challenge you've solved?"
- "Tell me about a complex problem you solved."
- "Tell me about a project you're proud of."
- "Tell me about a technical trade-off you made."

### Skeleton (this is what I look at)

```
S (~45 words)
  - ERP, ~5 clients; no reconciliation → users compared screen by screen
  - soccer club canteen: match days → spike of small transactions → 1h+

T (~25 words)
  - I designed and built bank reconciliation end to end, on my own, ~2 months

A (~150 words)
  1. Ingest: many sources (OFX, PDF, Open Finance via Pluggy)
     → one internal format; duplicates blocked by the DB, per bank account
  2. Rank top 5 open bills or invoices (date + amount) → user confirms
     WHY: silent wrong match >> few seconds of review
     only predictable cases (fees) automated via user rules
  3. Hardest part: many-to-many, WITH A LINK TABLE
     → warning when amounts don't add up
  (+1 sentence) sandbox → dogfooded with our finance team

R (~50 words)
  - match days: 1h+ → 10–15 min (client's feedback)
  - no incorrect reconciliations in production; still in use
  - most of the time, the top suggestion was the right one
```

If the story runs over 2:10, cut the testing sentence first.

### Full text (~300 words, reference only)

> At the time, I was working on an ERP used by about five clients, and it
> had no bank reconciliation. Users had to compare their bank statement
> with the system, line by line, to make sure nothing was missing. One
> client is a soccer club, and its canteen has lots of small
> transactions. On match days, reconciling everything could take over an
> hour.
>
> My job was to fix that. I designed and built the whole feature, end to
> end, by myself, in about two months, including testing.
>
> What I did first was the import side. Transactions come from different
> sources — OFX files, PDF statements, and bank APIs through Open
> Finance — so I normalized all of them into a single internal format.
> Duplicates are blocked by a unique index in the database, scoped to
> each bank account.
>
> Then came the key trade-off: automatic matching or suggestions. For
> each transaction, I rank the five closest open bills or invoices by
> date and amount, and the user just confirms. I chose suggestions
> because a wrong automatic match is a silent error, while reviewing
> takes a couple of seconds. I only automated the predictable cases,
> like bank fees, through user-defined rules.
>
> The hardest part was that matches aren't one-to-one. One payment can
> cover several bills, and one bill can be paid in parts. So I modeled
> it as many-to-many, with a link table and a warning when the amounts
> don't add up.
>
> Before rolling it out, I tested with sandbox data and then dogfooded
> it with our own finance team.
>
> As a result, the soccer club told me match days went from over an hour
> to about ten to fifteen minutes. It's still in production, no
> incorrect reconciliations have been reported, and most of the time
> users just confirm the top suggestion.

### Follow-ups

**"Why didn't you fully automate the matching?"**
> Because the risk of an automatic match being wrong — and nobody ever
> noticing — was far more expensive than the few seconds it takes a user
> to review and confirm. I only fully automated the predictable cases,
> like bank fees, through user-defined rules.

**"How did the ranking work?"**
> Two filters first: only candidates within 10 days of the transaction,
> and only ones not reconciled yet. Then a multi-key sort: exact amount
> and date first, then the closest date, then the smallest amount
> difference. The top five are shown as suggestions.

**"What if none of the five suggestions is right?"**
> The user searches manually, by payee name or description. If the bill
> doesn't exist yet, they can create it from the transaction, and the
> system pre-fills everything it can from the bank data.

**"How did you get the bank data?"**
> Through Pluggy, an Open Finance aggregator, plus OFX and PDF uploads.
> Two banks that aren't covered by Open Finance, Asaas and Simples, were
> integrated directly through their own APIs.

**"How did you prevent duplicate transactions on re-import?"**
> Each transaction gets a deduplication key built from the source's own
> ID, prefixed by the source: the FITID for OFX, Pluggy's transaction ID,
> and for PDFs, which have no ID, a hash of the normalized fields. A
> partial unique index on (bank account, key) enforces it for active
> rows, since we use soft delete. A unique violation counts as "already
> imported", not as an error. There's also a second check by a
> fingerprint of date, type, amount and description, and possible
> duplicates are flagged for the user to confirm.
>
> *Note: re-read this part of the code myself before using these details
> in an interview.*

**"Why per bank account, and not per bank?"**
> Because the OFX spec only guarantees that a FITID is unique within an
> account, not across the whole bank.

**"How did you handle amounts that didn't add up?"**
> Matches are stored in a link table, so one transaction can point to
> several bills and the other way around. When the sums don't match, the
> system warns the user and asks if they want to correct it.

**"How did you measure the result?"**
> From the client's feedback: match days went from over an hour to about
> ten to fifteen minutes. I didn't instrument it, which is part of what
> I'd do differently.

**"What would you do differently?"**
> I'd log from day one which suggestion users actually picked, and use
> that data to replace the cascading sort with a weighted score. Right
> now, a candidate with a closer date beats one with the exact amount,
> even when the exact amount is the more likely match.

### Pronunciation

| Word | Sounds like |
|---|---|
| reconciliation | rek-un-sil-ee-AY-shun |
| reconciling | rek-un-SY-ling |
| canteen | kan-TEEN |
| statement | STAYT-ment |
| scoped | skohpt (one syllable) |
| dogfooded | DOG-foo-did |
| suggestion | sug-JES-chun |

### Practice log

| Date | Time | Missed / stuck on |
|---|---|---|
| 2026-10-20 | 2:06 | Stuck on "payables/receivables" (switched to "bills or invoices"); left out "link table" |
