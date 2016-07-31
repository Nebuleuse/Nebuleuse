const UserMaskBase = 1 << 0;
const UserMaskOnlyId = 1 << 1;
const UserMaskAchievements = 1 << 2;
const UserMaskStats = 1 << 3;
const UserMaskAll = UserMaskStats | UserMaskAchievements;

const APIURL = "http://127.0.0.1:8000";