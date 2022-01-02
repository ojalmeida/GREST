<html lang="en-US">

    <link rel="stylesheet" href="css/style.css">
    <script src="/frontend/javascript/script.js"></script>

    <body>

        <header>


            <h1>OVERVIEW</h1>


        </header>

        <div id="wrapper">

            <div id="server-images">

                <div id="api" onmouseover="showAPIPane()">

                    <img src="/frontend/static/images/isometric-server-db.webp" alt="">
                    API

                </div>
                <svg id="arrow" width="546" height="16" viewBox="0 0 546 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M0.292893 7.29289C-0.0976311 7.68342 -0.0976311 8.31658 0.292893 8.70711L6.65685 15.0711C7.04738 15.4616 7.68054 15.4616 8.07107 15.0711C8.46159 14.6805 8.46159 14.0474 8.07107 13.6569L2.41421 8L8.07107 2.34315C8.46159 1.95262 8.46159 1.31946 8.07107 0.928932C7.68054 0.538408 7.04738 0.538408 6.65685 0.928932L0.292893 7.29289ZM545.707 8.70711C546.098 8.31658 546.098 7.68342 545.707 7.29289L539.343 0.928932C538.953 0.538408 538.319 0.538408 537.929 0.928932C537.538 1.31946 537.538 1.95262 537.929 2.34315L543.586 8L537.929 13.6569C537.538 14.0474 537.538 14.6805 537.929 15.0711C538.319 15.4616 538.953 15.4616 539.343 15.0711L545.707 8.70711ZM1 9L545 9V7L1 7L1 9Z" fill="#F3F3F3"/>
                </svg>

                <div id="db">

                    <img src="/frontend/static/images/isometric-server-api.webp" alt="">
                    DB

                </div>


            </div>

            <div id="api-pane" is-visible="false" onmouseleave="hideAPIPane(event)">

                <div id="header">API</div>
                <div id="navbar">

                    <div id="rate-limit-item" is-selected="false">
                        <svg width="48" height="48" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M10 18.4H16V38H10V18.4ZM21.2 10H26.8V38H21.2V10V10ZM32.4 26H38V38H32.4V26Z" fill="black"/>
                        </svg>
                        Rate <br>
                        Limit
                    </div>
                    <div id="path-mappings-item" is-selected="false">
                        <svg width="36" height="36" viewBox="0 0 36 36" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M17.9476 3C9.69166 3 2.99121 9.72 2.99121 18C2.99121 26.28 9.69166 33 17.9476 33C26.2035 33 32.9039 26.28 32.9039 18C32.9039 9.72 26.2035 3 17.9476 3ZM5.98248 18C5.98248 17.085 6.10214 16.185 6.29657 15.33L13.4457 22.5V24C13.4457 25.65 14.7918 27 16.437 27V29.895C10.5591 29.145 5.98248 24.105 5.98248 18ZM26.7569 26.1C26.368 24.885 25.2612 24 23.9152 24H22.4195V19.5C22.4195 18.675 21.7465 18 20.9239 18H11.9501V15H14.9413C15.7639 15 16.437 14.325 16.437 13.5V10.5H19.4283C21.0735 10.5 22.4195 9.15 22.4195 7.5V6.885C26.8017 8.655 29.9127 12.975 29.9127 18C29.9127 21.12 28.7012 23.97 26.7569 26.1Z" fill="black"/>
                        </svg>
                        Path <br>
                        Mappings
                    </div>
                    <div id="key-mappings-item" is-selected="false">
                        <svg width="36" height="36" viewBox="0 0 36 36" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M30 4.5H7.5C5.85 4.5 4.5 5.85 4.5 7.5V28.5C4.5 30.15 5.85 31.5 7.5 31.5H30C31.65 31.5 33 30.15 33 28.5V7.5C33 5.85 31.65 4.5 30 4.5ZM30 7.5V12H7.5V7.5H30ZM22.5 28.5H15V15H22.5V28.5ZM7.5 15H12V28.5H7.5V15ZM25.5 28.5V15H30V28.5H25.5Z" fill="black"/>
                        </svg>
                        Key <br>
                        Mappings
                    </div>
                    <div id="behaviors-item" is-selected="true">
                        <svg width="36" height="36" viewBox="0 0 36 36" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M14.67 16.74L12.54 18.87C11.52 17.835 10.53 16.5 9.855 14.46L12.765 13.725C13.245 15.06 13.92 15.975 14.67 16.74ZM16.5 9L10.5 3L4.5 9H9.03C9.06 10.215 9.15 11.31 9.315 12.255L12.225 11.52C12.12 10.8 12.045 9.945 12.03 9H16.5ZM31.5 9L25.5 3L19.5 9H23.985C23.835 14.52 22.065 16.125 20.175 17.82C19.425 18.48 18.66 19.2 18 20.145C17.49 19.41 16.905 18.825 16.305 18.285L14.19 20.4C15.585 21.675 16.5 22.71 16.5 25.5V33H19.5V25.5C19.5 22.47 20.565 21.51 22.185 20.055C24.255 18.195 26.805 15.885 26.985 9H31.5V9Z" fill="black"/>
                        </svg>
                        Behaviors
                    </div>


                </div>

                <div id="rate-limit-menu"></div>
                <div id="path-mappings-menu"></div>
                <div id="key-mappings-menu"></div>
                <div id="behaviors-menu"></div>


            </div>

            <div id="db-pane" is-visible="false"></div>

        </div>

    </body>

</html>
