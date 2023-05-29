package unionpay

const (
	kFrontTransTemplate = `
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" /></head>
<body onload="onLoad();">
<form id="unionpay" action="{{.Action}}" method="POST">
{{range $k, $v := .Values}}
<input type="hidden" name="{{$k}}" id="{{$k}}" value="{{index $v 0}}" />
{{end}}
</form>
<script type="text/javascript">
<!--
function onLoad()
{
document.getElementById("unionpay").submit();
}
//-->
</script>
</body>
</html>
`
)
