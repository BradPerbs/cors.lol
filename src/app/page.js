'use client'
import Image from "next/image";
import Header from "@/app/components/Header";
import bannerBgImg from './assets/images/banner-bg-img.png';
import Accordion from 'react-bootstrap/Accordion';
import Footer from "@/app/components/Footer";
export default function Home() {
  return (
      <div className='home-page'>
          <Header/>
          <div className="home-banner" style={{backgroundImage: `url(${bannerBgImg.src})`}}>
              <div className="container">
                  <div className="row">
                      <div className="col-md-6 align-self-center">
                          <div className="banner-content">
                              <p className="large cyan">Tired of CORS errors?</p>
                              <h1>Super easy, free to use CORS proxy</h1>
                              <p className="large">Stop fighting with CORS errors, simply proxy your requests through cors.lol, an easy to use, open-source CORS proxy.</p>
                              <div className="banner-btns">
                                  <a href="#getStarted" className="btn-style colored-border">Get Started <Image src={require('./assets/images/icons/right-arrow-icon.png')} alt='right-arrow-icon' className='right-arrow-icon' /> </a>
                                  <a href="#"><Image src={require('./assets/images/icons/github-icon.png')} alt='github-icon' className='github-icon' /></a>
                              </div>
                          </div>
                      </div>
                      <div className="col-md-6 align-self-center">
                          <div className="banner-img-wrapper">
                              <Image src={require('./assets/images/headerimg.svg')} alt='banner-img' className='banner-img'/>
                          </div>
                      </div>
                  </div>
              </div>
          </div>
          <div className="how-do-use-it" id="getStarted" name="getStarted">
              <div className="container">
                  <div className="section-header">
                      <p className="cyan medium">Get Started</p>
                      <h2>How do I use it?</h2>
                      <p>Simply add the proxy url before the url you want to call, like the following example</p>
                  </div>
                  <form action="#" >
                      <div className="form-group">
                          <div className="url">
                              <Image src={require('./assets/images/icons/globe-icon.png')} alt='github-icon' className='globe-icon'/>
                              <p className="large white">https://api.cors.lol/</p>
                          </div>
                          <input type="text" className="form-control" placeholder='yoururl.com/yourapiendpoint' />
                      </div>
                      <a href="#" className="btn-style colored-border">Copy to clipboard</a>
                  </form>
              </div>
          </div>
          <div className="flexible-plan" id="pricing" name="pricing">
              <div className="container">
                  <div className="section-header">
                      <p className="cyan">Flexible Plans</p>
                      <h2>Choose the right fit for <br/> your usage</h2>
                  </div>
                  <div className="row">
                      <div className="col-md-6">
                          <div className="plan-main-wrapper">
                              <span className="popular-label">POPULAR</span>
                              <p className="large white">BASIC</p>
                              <span>Unlimited free of use for non commercial use</span>
                              <h3 className='large'>Free</h3>
                              <ul className="features">
                                  <li className='active'>
                                      <div className="list">
                                          <Image src={require('./assets/images/icons/tick-icon.png')} alt='github-icon'
                                                 className='feature-icon'/>
                                          <span>Rate Limited</span></div>
                                  </li>
                                  <li className='active'>
                                      <div className="list">
                                          <Image src={require('./assets/images/icons/tick-icon.png')} alt='github-icon'
                                                 className='feature-icon'/>
                                          <span>Unlimited requests</span></div>
                                  </li>
                                  <li>
                                      <div className="list">
                                          <Image src={require('./assets/images/icons/close-icon.png')} alt='github-icon'
                                                 className='feature-icon'/>
                                          <span>Support</span></div>
                                  </li>
                                  <li>
                                      <div className="list">
                                          <Image src={require('./assets/images/icons/close-icon.png')} alt='github-icon'
                                                 className='feature-icon'/>
                                          <span>Up to 10 MB per request</span></div>
                                  </li>
                              </ul>
                              <a href="#" className="btn-style feature-btn dark">Get Started</a>
                          </div>
                      </div>
                      <div className="col-md-6">
                          <div className="plan-main-wrapper popular">
                              <span className="popular-label">POPULAR</span>
                              <p className="large white">LIFETIME PRO</p>
                              <span>Unlimited free of use for non commercial use</span>
                              <h3 className='large'>$79 <span>/ LIFETIME</span></h3>
                              <ul className="features">
                                  <li className='active'>
                                      <div className="list">
                                          <Image src={require('./assets/images/icons/tick-icon.png')} alt='github-icon' className='feature-icon'/>
                                          <span>No Rate Limit</span></div>
                                  </li>
                                  <li className='active'>
                                      <div className="list">
                                          <Image src={require('./assets/images/icons/tick-icon.png')} alt='github-icon'
                                                 className='feature-icon'/>
                                          <span>Unlimited Requests</span></div>
                                  </li>
                                  <li className='active'>
                                      <div className="list">
                                          <Image src={require('./assets/images/icons/tick-icon.png')} alt='github-icon'
                                                 className='feature-icon'/>
                                          <span>Technical Support</span></div>
                                  </li>
                                  <li className='active'>
                                      <div className="list">
                                          <Image src={require('./assets/images/icons/tick-icon.png')} alt='github-icon'
                                                 className='feature-icon'/>
                                          <span>Up to 100 MB per request</span></div>
                                  </li>
                              </ul>
                              <a href="#" className="btn-style feature-btn colored-border">Get Started</a>
                          </div>
                      </div>
                  </div>
              </div>
          </div>
          <div className="our-blog">
              <div className="container">
                  <div className="section-heading">
                      <p className="cyan">Our Blog</p>
                      <h2 className="large">Take a look at the latest <br/> articles from</h2>
                  </div>
                  <div className="row">
                      <div className="col-md-6 col-lg-4">
                          <div className="blog-main-wrapper">
                              <div className="img-wrapper">
                                  <Image src={require('./assets/images/blog-img-1.png')} alt='blog-img'
                                         className='blog-img'/>
                              </div>
                              <div className="content-wrapper">
                                  <ul className="list">
                                      <li>saas</li>
                                      <li>development</li>
                                  </ul>
                                  <h3 className="large">Unlocking the potential of SaaS: A journey to business
                                      success</h3>
                                  <div className="footer-text">
                                      <ul>
                                          <li>July 17, 2023</li>
                                          <li>3 min read</li>
                                      </ul>
                                  </div>
                              </div>
                          </div>
                      </div>
                      <div className="col-md-6 col-lg-4">
                          <div className="blog-main-wrapper">
                              <div className="img-wrapper">
                                  <Image src={require('./assets/images/blog-img-2.png')} alt='blog-img' className='blog-img'/>
                              </div>
                              <div className="content-wrapper">
                                  <ul className="list">
                                      <li>nocode</li>
                                      <li>tools</li>
                                  </ul>
                                  <h3 className="large">Webflow vs Framer - Which one is best for you and your business.</h3>
                                  <div className="footer-text">
                                      <ul>
                                          <li>July 12, 2023</li>
                                          <li>4 min read</li>
                                      </ul>
                                  </div>
                              </div>
                          </div>
                      </div>
                      <div className="col-md-6 col-lg-4">
                          <div className="blog-main-wrapper">
                              <div className="img-wrapper">
                                  <Image src={require('./assets/images/blog-img-3.png')} alt='blog-img'
                                         className='blog-img'/>
                              </div>
                              <div className="content-wrapper">
                                  <ul className="list">
                                      <li>organization</li>
                                  </ul>
                                  <h3 className="large">The 3 Best Social Media Analytics Tools For Competitor Analysis</h3>
                                  <div className="footer-text">
                                      <ul>
                                          <li>July 19, 2023</li>
                                          <li>4 min read</li>
                                      </ul>
                                  </div>
                              </div>
                          </div>
                      </div>
                  </div>
              </div>
          </div>
          <div className="faqs">
              <div className="container">
                  <div className="row">
                      <div className="col-lg-4 col-md-5">
                          <div className="content-wrapper">
                              <p className="cyan">Frequently Asked Questions</p>
                              <h2>Your CORS
                                  solutions guide</h2>
                              <p>We understand that you may have some queries before starting using corslol. Here you
                                  can find some answered questions.</p>
                          </div>
                      </div>
                      <div className="col-lg-8 col-md-7">
                          <Accordion defaultActiveKey="0">
                              <Accordion.Item eventKey="0">
                                  <Accordion.Header>What is a CORS proxy?</Accordion.Header>
                                  <Accordion.Body>
                                      <p>A CORS proxy service allows web developers to bypass the same-origin policy enforced by web browsers. This service acts as an intermediary between your web application and the target server, enabling your application to make cross-origin requests without encountering CORS (Cross-Origin Resource Sharing) errors.</p>
                                  </Accordion.Body>
                              </Accordion.Item>
                              <Accordion.Item eventKey="1">
                                  <Accordion.Header>Is this CORS proxy service secure to use?</Accordion.Header>
                                  <Accordion.Body>
                                      <p>While our CORS proxy service is designed to be secure, it is important to understand that using any proxy introduces potential security risks. We recommend using the proxy only for development purposes and not for sensitive or production data. Always review the code and understand the security implications before deploying it in your environment.</p>
                                  </Accordion.Body>
                              </Accordion.Item>
                              <Accordion.Item eventKey="2">
                                  <Accordion.Header>How can I set up and use this CORS proxy service?</Accordion.Header>
                                  <Accordion.Body>
                                      <p>Using our CORS proxy service is simple. You don't need to install anything. Just prepend your target API URL with https://api.cors.lol/. For example, if you want to access https://example.com/api, you would use https://api.cors.lol/https://example.com/api. This way, all your cross-origin requests will be routed through our proxy, avoiding CORS restrictions.</p>
                                  </Accordion.Body>
                              </Accordion.Item>
                              <Accordion.Item eventKey="3">
                              <Accordion.Header>How secure is my data on a SaaS website?</Accordion.Header>
                              <Accordion.Body>
                                  <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
                                      eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad
                                      minim veniam, quis nostrud exercitation ullamco laboris nisi ut
                                      aliquip ex ea commodo consequat. Duis aute irure dolor in
                                      reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
                                      pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
                                      culpa qui officia deserunt mollit anim id est laborum.</p>
                              </Accordion.Body>
                          </Accordion.Item>

                          </Accordion>
                      </div>
                  </div>
              </div>
          </div>
          <Footer/>
      </div>
  );
}
