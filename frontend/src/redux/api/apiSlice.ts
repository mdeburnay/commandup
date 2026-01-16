import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { logout } from "../slices/authSlice";

const baseQuery = fetchBaseQuery({
	baseUrl: "http://localhost:8080/api",
	prepareHeaders: (headers, { getState }) => {
		const token = (getState() as any).auth.accessToken;
		if (token) {
			headers.set("authorization", `Bearer ${token}`);
		}
		return headers;
	},
});

const baseQueryWithReauth = async (args: any, api: any, extraOptions: any) => {
	let result = await baseQuery(args, api, extraOptions);

	if (result.error && result.error.status === 401) {
		api.dispatch(logout());
	}
	return result;
};

export const apiSlice = createApi({
	baseQuery: baseQueryWithReauth,
	endpoints: (builder) => ({
		login: builder.mutation({
			query: (credentials) => ({
				url: "/auth/login",
				method: "POST",
				body: { ...credentials },
			}),
		}),
		signup: builder.mutation({
			query: (userData) => ({
				url: "/auth/signup",
				method: "POST",
				body: { ...userData },
			}),
		}),
	}),
});

export const { useLoginMutation, useSignupMutation } = apiSlice;
