import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { logout, setCredentials } from "../slices/authSlice";

const baseQuery = fetchBaseQuery({
	baseUrl: "http://localhost:8080/api",
	credentials: "include", // Include cookies in requests
});

const baseQueryWithReauth = async (args: any, api: any, extraOptions: any) => {
	let result = await baseQuery(args, api, extraOptions);

	// If we get a 401, try to refresh the token
	if (result.error && result.error.status === 401) {
		const refreshResult = await baseQuery(
			{ url: "/auth/refresh", method: "POST" },
			api,
			extraOptions
		);

		if (refreshResult.data) {
			// Token refreshed, retry the original request
			result = await baseQuery(args, api, extraOptions);
		} else {
			// Refresh failed, logout user
			api.dispatch(logout());
		}
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
			async onQueryStarted(arg, { dispatch, queryFulfilled }) {
				try {
					const { data } = await queryFulfilled;
					dispatch(
						setCredentials({
							user: data.user,
						})
					);
				} catch (err) {
					// TODO: Handle error
				}
			},
		}),
		signup: builder.mutation({
			query: (userData) => ({
				url: "/auth/signup",
				method: "POST",
				body: { ...userData },
			}),
			async onQueryStarted(arg, { dispatch, queryFulfilled }) {
				try {
					const { data } = await queryFulfilled;
					dispatch(
						setCredentials({
							user: data.user,
						})
					);
				} catch (err) {
					// TODO: Handle error
				}
			},
		}),
		logout: builder.mutation({
			query: () => ({
				url: "/auth/logout",
				method: "POST",
			}),
			async onQueryStarted(arg, { dispatch, queryFulfilled }) {
				try {
					await queryFulfilled;
					dispatch(logout());
				} catch (err) {
					// TODO: Handle error
					dispatch(logout());
				}
			},
		}),
	}),
});

export const { useLoginMutation, useSignupMutation, useLogoutMutation } =
	apiSlice;
