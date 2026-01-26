import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface User {
	id: string;
	email: string;
	username: string;
}

interface AuthState {
	user: User | null;
}

const initialState: AuthState = {
	user: localStorage.getItem("user")
		? JSON.parse(localStorage.getItem("user")!)
		: null,
};

const authSlice = createSlice({
	name: "auth",
	initialState,
	reducers: {
		setCredentials: (
			state,
			action: PayloadAction<{
				user: User;
			}>
		) => {
			const { user } = action.payload;
			state.user = user;
			localStorage.setItem("user", JSON.stringify(user));
		},
		logout: (state) => {
			state.user = null;
			localStorage.removeItem("user");
		},
	},
});

export const { setCredentials, logout } = authSlice.actions;

export default authSlice.reducer;

export const selectCurrentUser = (state: { auth: AuthState }) =>
	state.auth.user;
