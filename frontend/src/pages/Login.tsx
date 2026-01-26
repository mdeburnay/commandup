// Packages
import { useState } from "react";
import { useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";

// Redux
import { useLoginMutation } from "../redux/api/apiSlice";
import { setCredentials } from "../redux/slices/authSlice";

// Components
import { Button } from "../components/Buttons/PrimaryButton";

export const Login = () => {
	const [email, setEmail] = useState<string>("");
	const [password, setPassword] = useState<string>("");
	const [error, setError] = useState<string>("");

	const dispatch = useDispatch();
	const navigate = useNavigate();
	const [login, { isLoading }] = useLoginMutation();

	const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		try {
			await login({ email, password }).unwrap();
			setError("");
			navigate("/");
		} catch (err: any) {
			setError(err?.data?.message || err.error || "Login failed");
		}
	};

	return (
		<>
			<form
				className="flex justify-center items-center flex-col"
				onSubmit={handleSubmit}
			>
				<input
					className="flex w-80 p-1 rounded-md m-2"
					type="text"
					placeholder="Email"
					value={email}
					onChange={(e) => setEmail(e.target.value)}
				/>
				<input
					className="flex w-80 p-1 rounded-md m-2"
					type="password"
					placeholder="Password"
					value={password}
					onChange={(e) => setPassword(e.target.value)}
				/>
				<button
					type="submit"
					disabled={isLoading}
					className="flex p-1 rounded-md m-2 hover:cursor-pointer disabled:opacity-50"
				>
					{isLoading ? "Logging in..." : "Login"}
				</button>
				{error && <div className="text-red-500">{error}</div>}
				<div className="p-4">Don't have an account? </div>
				<Button text="Sign Up" url="/signup" />
			</form>
		</>
	);
};
