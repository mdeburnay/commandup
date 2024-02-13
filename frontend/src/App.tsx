// Dependencies
import {
  BrowserRouter as Router,
  Route,
  Routes,
  useNavigate,
} from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

// Styles
import "./index.css";

// Components
import { Container } from "./components/Container";

// Pages
import { Home } from "./pages/Home";
import { Login } from "./pages/Login";

// Create a client
const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <Routes>
          <Route path="/" element={<Container />}>
            <Route index element={<Home />} />
            <Route index element={<Login />} />
          </Route>
        </Routes>
      </Router>
    </QueryClientProvider>
  );
}

export default App;
