package cc.cryptopunks.portal.onboarding

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.material.Button
import androidx.compose.material.Divider
import androidx.compose.material.MaterialTheme
import androidx.compose.material.OutlinedTextField
import androidx.compose.material.Surface
import androidx.compose.material.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import cc.cryptopunks.portal.compose.AstralTheme
import cc.cryptopunks.portal.core.mobile.Core
import org.koin.compose.koinInject

@Composable
fun OnBoardingScreen(core: Core = koinInject()) = OnboardingScreen(core::setup)

@Composable
fun OnboardingScreen(
    proceed: (String) -> Unit
) = Column(
    modifier = Modifier
        .fillMaxSize()
        .padding(32.dp),
    horizontalAlignment = Alignment.CenterHorizontally,
    verticalArrangement = Arrangement.Center,
) {

    var alias by remember { mutableStateOf("") }
    var error by remember { mutableStateOf(false) }
    OutlinedTextField(
        value = alias,
        modifier = Modifier.fillMaxWidth(),
        onValueChange = {
            alias = it
            if (alias != "") error = false
        },
        label = { Text(text = "user alias") },
        isError = error
    )
    Spacer(modifier = Modifier.size(8.dp))
    Button(
        modifier = Modifier.fillMaxWidth(),
        onClick = {
            if (alias == "") error = true
            else proceed(alias)
        }
    ) {
        Text(text = "Create new user")
    }

    Row(
        verticalAlignment = Alignment.CenterVertically
    ) {
        Divider(
            modifier = Modifier
                .height(2.dp)
                .weight(1f)
        )
        Text(
            modifier = Modifier.padding(24.dp),
            text = "or",
            color = Color.White,
            style = MaterialTheme.typography.subtitle1
        )
        Divider(
            modifier = Modifier
                .height(1.dp)
                .weight(1f)
        )
    }

    Button(
        modifier = Modifier.fillMaxWidth(),
        onClick = { proceed("") }
    ) {
        Text(text = "Confirm already claimed")
    }
}

@Preview
@Composable
fun OnboardingScreenPreview() = AstralTheme {
    Surface {
        OnboardingScreen {}
    }
}