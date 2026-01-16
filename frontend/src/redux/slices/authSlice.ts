import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface User {
	id: string;
	email: string;
	username: string;
}

interface AuthState {
	user: User | null;
	accessToken: string | null;
	refreshToken: string | null;
}

const initialState: AuthState = {
	user: localStorage.getItem("user")
		? JSON.parse(localStorage.getItem("user")!)
		: null,
	accessToken: localStorage.getItem("accessToken"),
	refreshToken: localStorage.getItem("refreshToken"),
};

const authSlice = createSlice({
	name: "auth",
	initialState,
	reducers: {
		setCredentials: (
			state,
			action: PayloadAction<{
				user: User;
				accessToken: string;
				refreshToken?: string;
			}>
		) => {
			const { user, accessToken, refreshToken } = action.payload;
			state.user = user;
			state.accessToken = accessToken;
			if (refreshToken) {
				state.refreshToken = refreshToken;
			}

			localStorage.setItem("accessToken", accessToken);
			if (refreshToken) {
				localStorage.setItem("refreshToken", refreshToken);
			}
			localStorage.setItem("user", JSON.stringify(user));
		},
		logout: (state) => {
			state.user = null;
			state.accessToken = null;
			state.refreshToken = null;
			localStorage.removeItem("accessToken");
			localStorage.removeItem("refreshToken");
			localStorage.removeItem("user");
		},
	},
});

export const { setCredentials, logout } = authSlice.actions;

export default authSlice.reducer;

export const selectCurrentUser = (state: { auth: AuthState }) =>
	state.auth.user;
export const selectCurrentToken = (state: { auth: AuthState }) =>
	state.auth.accessToken;
