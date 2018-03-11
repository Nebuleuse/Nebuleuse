const UserMaskBase = 1 << 0;
const UserMaskOnlyId = 1 << 1;
const UserMaskAchievements = 1 << 2;
const UserMaskStats = 1 << 3;
const UserMaskAll = UserMaskStats | UserMaskAchievements;

const APIURL = '//' + window.location.host;
