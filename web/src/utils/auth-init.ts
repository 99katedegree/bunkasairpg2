import Cookies from "js-cookie";

export const getAuthInit = (): RequestInit => ({
  headers: { Authorization: `Bearer ${Cookies.get("authToken")}` },
});
