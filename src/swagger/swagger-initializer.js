window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">
  window.ui = SwaggerUIBundle({
    urls: [
      { url: "/openapi/shortn.yml", name: "Shortn API" },
      { url: "/openapi/blowup.yml", name: "Blowup API" },
    ],
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    plugins: [SwaggerUIBundle.plugins.DownloadUrl],
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};
