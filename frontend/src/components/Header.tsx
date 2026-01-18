// Packages
import { useSelector, useDispatch } from "react-redux";
import { logout, selectCurrentUser } from "../redux/slices/authSlice";

// Components
import { Button } from "./Buttons/PrimaryButton";
import { FileUpload } from "./FileUpload";

export const Header = () => {
	const user = useSelector(selectCurrentUser);
	const dispatch = useDispatch();

	const handleLogout = () => {
		dispatch(logout());
	};

	return (
		<header className="flex justify-end px-12">
			<FileUpload />
			{user ? (
				<Button text="Logout" onClick={handleLogout} />
			) : (
				<Button text="Login" url="/login" />
			)}
		</header>
	);
};
