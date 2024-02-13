// Components
import { Button } from "./Button";
import { FileUpload } from "./FileUpload";

export const Header = () => {
  return (
    <header className="flex justify-end px-12">
      <FileUpload />
      <Button text="Login" />
    </header>
  );
};
