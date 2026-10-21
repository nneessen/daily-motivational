# Morning texts

Every weekday around 5:20 AM Denver time, GitHub runs this program. It asks Claude for 4 options for the morning text to the team (3 short, 1 longer) and emails them to you (the `TO_ADDRESS` secret). Pick one, tweak it, send it.

## Change what it writes
- **Tone and rules:** edit `voice.md`. It's the whole prompt, with your real example message in it.
- **This week's news:** write a line or two in `focus.txt` (like "Contest ends Friday."). It gets worked into the next texts. Clear it out when it's old. The repo is public, so only write what you'd be fine with anyone reading.
- **Angles:** the list of ideas the texts rotate through is `angles` in `prompt.go`.

Edits on github.com count too. The next morning's run picks them up.

## One-time setup
1. **Keys:** both are the ones commissionTracker already uses.
   - The Anthropic key costs about $1–3 a month here.
   - Mailgun sends one email a day from `morning@updates.thestandardhq.com`.
2. **Local testing:** `cp .env.example .env` and fill in the two keys and `TO_ADDRESS`.
3. **GitHub secrets:** `gh secret set -f .env` loads all four from `.env`: `ANTHROPIC_API_KEY`, `MAILGUN_API_KEY`, `MAILGUN_DOMAIN` and `TO_ADDRESS`.

## Run it
```
scripts/test_run.sh          # checks + prints today's email, sends nothing
scripts/test_run.sh --send   # checks + sends one real email
gh workflow run daily.yml    # run the GitHub job now
```

## Schedule
GitHub cron only runs in UTC, so `.github/workflows/daily.yml` has two lines. One covers daylight time and one covers standard time, and the job skips whichever doesn't match Denver's clock that day. GitHub sometimes starts scheduled jobs 5–30 minutes late, which is why it's set for 5:17.

The repo is public because public repos get unlimited free Actions minutes, and theStandardHQ uses up this account's private-repo minutes every month. The catch: GitHub pauses schedules in a public repo after 60 days with no commits. It emails a warning first, and any edit (even to `focus.txt`) resets the clock. The keys live in GitHub secrets, which stay hidden in a public repo.
