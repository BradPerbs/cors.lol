import Image from "next/image";

export default function Footer(){
    return (
        <footer className="footer">
            <div className="container">
                <div className="row">
                    <div className="col-md-4">
                        <div className="footer-content">
                            <div className="logo-wrapper">
                                <a href="#">
                                    <Image src={require('../assets/images/icons/logo.png')} alt='logo'/>
                                </a>
                            </div>
                            <div className="text-wrapper">
                                <span>Open source and free to use, our CORS proxy service allows seamless  cross-origin requests. Enhance your web applications functionality with  secure, efficient, and reliable API access.</span>
                            </div>
                            <ul className="social-icons">
                                <li><a href="#"> <Image src={require('../assets/images/icons/x-iocn.png')} alt='logo'/></a>
                                </li>
                                <li><a href="#"> <Image src={require('../assets/images/icons/github-icon.png')}
                                                        alt='logo'/></a></li>
                                <li><a href="#"> <Image src={require('../assets/images/icons/circles-icon.png')}
                                                        alt='logo'/></a></li>
                            </ul>
                            <span className="copy-right-text">Copyright ©2024 Cors.lol</span>
                        </div>
                    </div>
                    <div className="col-md-3 col-4">
                        <div className="footer-links">
                            <strong className="white">Pages</strong>
                            <ul className="links">
                                <li><a href="#">Home</a></li>
                                <li><a href="#">How to use it</a></li>
                                <li><a href="#">Pricing</a></li>
                                <li><a href="#">Blog</a></li>
                            </ul>
                        </div>
                    </div>
                    <div className="col-md-3 col-4">
                        <div className="footer-links">
                            <strong className="white">Pages</strong>
                            <ul className="links">
                                <li><a href="#">Integration</a></li>
                                <li><a href="#">Blog</a></li>
                                <li><a href="#">Blog Article</a></li>
                                <li><a href="#">Contact</a></li>
                            </ul>
                        </div>
                    </div>
                    <div className="col-md-2 col-4">
                        <div className="footer-links">
                            <strong className="white">Account</strong>
                            <ul className="links">
                                <li><a href="#">Login</a></li>
                                <li><a href="#">Register</a></li>
                                <li><a href="#">Privacy Policy</a></li>
                                <li><a href="#">Support</a></li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        </footer>
    )
}