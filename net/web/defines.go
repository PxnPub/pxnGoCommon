package web;



const LogPrefixWeb   = "[Web] ";
const DefaultBindWeb = "tcp://127.0.0.1:8000";



// mimes
const(
	Mime_Text = "text/plain"
	Mime_HTML = "text/html"
	Mime_JSON = "application/json"
	Mime_SVG  = "image/svg+xml"
);



// html tags
const(
	Tag_Title          = "Title"
	Tag_Page           = "Page"
	Tag_FavIcon        = "FavIcon"
	Tag_FilesCSS       = "FilesCSS"
	Tag_RawCSS         = "RawCSS"
	// bootstrap
	Tag_WithBootstrap  = "WithBootstrap"
	Tag_WithBootsIcons = "WithBootstrapIcons"
	Tag_WithTooltips   = "WithTooltips"
	Tag_WithFloatingUI = "WithBootstrapPopper"
	// jquery
	Tag_WithJQuery     = "WithJQuery"
	// datatables
	Tag_WithDataTables = "WithDataTables"
	// echarts
	Tag_WithECharts    = "WithECharts"
	Tag_AppendHead     = "AppendHead"
	Tag_AppendHeader   = "AppendHeader"
	Tag_AppendFooter   = "AppendFooter"
);

// include urls
const(
	// bootstrap
	URL_BootstrapCSS           = "https://cdn.jsdelivr.net/npm/bootstrap@{{VERSION}}/dist/css/bootstrap.min.css";
	URL_BootstrapJS            = "https://cdn.jsdelivr.net/npm/bootstrap@{{VERSION}}/dist/js/bootstrap.bundle.min.js";
	URL_BootsIconsCSS          = "https://cdn.jsdelivr.net/npm/bootstrap-icons@1.13.1/font/bootstrap-icons.min.css";
	URL_FloatingJS             = "https://cdn.jsdelivr.net/npm/@popperjs/core@2.11.8/dist/umd/popper.min.js";
	// jquery
	URL_JQueryJS               = "https://cdnjs.cloudflare.com/ajax/libs/jquery/3.7.1/jquery.min.js";
	// datatables
	URL_DataTablesJS           = "https://cdn.datatables.net/2.3.1/js/dataTables.min.js";
	URL_DataTablesBootstrapJS  = "https://cdn.datatables.net/2.3.1/js/dataTables.bootstrap5.min.js";
	URL_DataTablesBootstrapCSS = "https://cdn.datatables.net/2.3.1/css/dataTables.bootstrap5.min.css";
	URL_DataTablesScrollerJS   = "https://cdn.datatables.net/scroller/2.4.3/js/dataTables.scroller.min.js";
	URL_DataTablesScrollerCSS  = "https://cdn.datatables.net/scroller/2.4.3/css/scroller.bootstrap5.min.css";
	URL_DataTablesPageResizeJS = "https://cdn.datatables.net/plug-ins/2.3.1/features/pageResize/dataTables.pageResize.min.js";
	// echarts
	URL_EChartsJS              = "https://cdnjs.cloudflare.com/ajax/libs/echarts/5.6.0/echarts.min.js"
);
