package cc.cryptopunks.portal.permission

import android.content.Intent
import android.net.Uri
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material.Button
import androidx.compose.material.MaterialTheme
import androidx.compose.material.Scaffold
import androidx.compose.material.Text
import androidx.compose.material.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import cc.cryptopunks.portal.Permissions
import cc.cryptopunks.portal.compose.AstralTheme

class PermissionActivity : ComponentActivity() {

    private val requestPermission = registerForActivityResult(
        ActivityResultContracts.RequestMultiplePermissions()
    ) { result ->
        val rejected = result.filterValues { !it }.keys.toTypedArray()
        val intent = Permissions.result(rejected)
        setResult(RESULT_OK, intent)
        finish()
    }

    private val requestStorage = registerForActivityResult(
        ActivityResultContracts.StartActivityForResult()
    ) {
        finish()
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        val message = Permissions.getMessage(intent)
        val required = Permissions.getRequired(intent).toMutableSet()

        required.isNotEmpty() || return
        setContent {
            AstralTheme {
                PermissionsView(message) {
                    val perm = required.find { perm ->
                        perm.startsWith("android.settings")
                    }
                    if (perm != null) {
                        val uri = Uri.parse("package:$packageName")
                        val intent = Intent(perm, uri)
                        requestStorage.launch(intent)
                    } else if (required.isNotEmpty()) {
                        requestPermission.launch(required.toTypedArray())
                    }
                }
            }
        }
    }
}

@Preview
@Composable
private fun PermissionsPreview() {
    AstralTheme {
        PermissionsView(message = "Test message") {

        }
    }
}

@Composable
private fun PermissionsView(
    message: String,
    onClick: () -> Unit,
) = Scaffold(
    topBar = {
        TopAppBar(
            title = {
                Text("Permission required")
            }
        )
    }
) { paddingValues ->
    Column(
        Modifier
            .fillMaxSize()
            .padding(paddingValues),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
    ) {
        // setup permissions rationale
        Text(
            text = message,
            modifier = Modifier.padding(64.dp),
            textAlign = TextAlign.Center,
            style = MaterialTheme.typography.h5,
        )
        // setup grant permissions button
        Button(
            content = {
                Text("grant permissions")
            },
            onClick = onClick,
        )
    }
}
