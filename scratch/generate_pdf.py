from fpdf import FPDF

class PDF(FPDF):
    pass

pdf = PDF()
pdf.add_page()
pdf.set_font("helvetica", size=12)
pdf.multi_cell(0, 10, """
MEETEXT MEETING NOTES
Date: May 20, 2026
Attendees: Sarah (Tech Lead), John (Product Manager), Acme Client Representatives.

John opened the meeting discussing mobile authentication options.
Sarah proposed using OAuth2 with JWT tokens:
- "I will implement the JWT authentication endpoint and token verification by Friday. We will use RS256 for signing."
- John agreed and said, "Let's use PostgreSQL instead of MongoDB to store credentials and session tables."

Blockers:
The payment gateway integration (Stripe API migration) is blocked because the client has not provided production API keys. This blocker is critical and must be resolved by next Tuesday.

Client Requests:
Acme Client requested a PDF export feature for meeting results by end of sprint.

Documentation tasks:
We need to document the database schema changes and update the API documentation.
""")
pdf.output("test_meeting.pdf")
