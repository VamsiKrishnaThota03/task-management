import { redirect } from "next/navigation";

// Kept so old links keep working; the sign-up form lives at /register.
export default function SignupPage() {
    redirect("/register");
}
