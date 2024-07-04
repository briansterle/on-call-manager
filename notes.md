a) App Overview:

Name: OnCallManager
Purpose: Coordinate emergency anointing calls among priests across multiple parishes
Platforms: Android and iOS

b) User Roles:

Priests: Receive calls, mark availability, respond to emergencies
Parish Secretaries: Can input emergency requests
Administrator: Manages priest schedules and overall system

c) Core Features:

Emergency Notification System:

Receive notifications via email
Convert email to in-app notification
Push notification to on-call priests


On-Call Calendar:

Display weekly/monthly view of priest schedules
Allow administrators to set and modify schedules


Call Response System:

Priests can accept a call within the app
Once accepted, notify other priests that the call is taken
Prevent multiple acceptances for the same call


Call History:

Log each emergency call
Record which priest responded
Provide simple reporting interface


User Management:

Add/remove priests and other users
Set user roles and permissions



d) Technical Requirements:

Real-time database for call status and priest availability
Push notification system for both Android and iOS
Email integration for receiving emergency requests
Secure authentication system
VOIP system integration for voicemail transcription (future feature)

e) Non-Functional Requirements:

User-friendly interface suitable for less tech-savvy users
Fast response time (< 2 seconds) for critical operations
High availability (99.9% uptime)
Data encryption for sensitive information


Wireframes:
For the wireframes, I'll describe key screens. In practice, these would be visual mockups, but I'll provide text descriptions:

a) Home Screen:

Current on-call status (prominent)
Quick action button to respond to active call
Upcoming on-call schedule (next 3-5 shifts)
Recent call history (last 3-5 calls)

b) Calendar View:

Monthly calendar grid
Color-coded shifts for each priest
Tap day to see detailed schedule

c) Active Call Screen:

Emergency details (location, patient info if available)
Large "Accept Call" button
Option to view more info or contact requester

d) Call History:

List view of past calls
Date, time, responding priest, and basic call info
Option to tap for more details

e) Profile/Settings:

Personal information
Notification preferences
Availability settings

f) Admin Panel:

User management section
Schedule management tools
System settings and configuration

